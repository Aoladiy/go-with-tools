package auth

import (
	"context"
	"errors"
	"fmt"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/config"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"
	"go-with-tools/internal/helpers"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessExp  = time.Minute * 15
	refreshExp = time.Hour * 24 * 7
)

type Service struct {
	q *queries.Queries
	p *pgxpool.Pool
	c config.Config
}

func New(q *queries.Queries, p *pgxpool.Pool, c config.Config) *Service {
	return &Service{q: q, p: p, c: c}
}

func (s *Service) SignUp(ctx context.Context, request DTO.SignUpRequest) (DTO.JWTResponse, *errs.AppError) {
	if len(request.Password) < 8 {
		return DTO.JWTResponse{}, errs.BadRequest(errors.New("password must be at least 8 characters"))
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}

	var jwtResponse DTO.JWTResponse
	appErr := helpers.WithTx(ctx, s.p, s.q, func(timeout context.Context, q *queries.Queries) *errs.AppError {
		adminUser, err := q.CreateAdminUser(timeout, queries.CreateAdminUserParams{
			Email:        request.Email,
			PasswordHash: string(password),
		})
		if err != nil {
			if pgErr, isUniqueViolation := errs.IsUniqueViolation(err); isUniqueViolation {
				return errs.UniqueViolation(err, pgErr)
			}
			return errs.Internal(err)
		}
		var appErr *errs.AppError
		jwtResponse, appErr = generateJWTResponse(s.c.JwtSecret, strconv.FormatInt(adminUser.ID, 10))
		if appErr != nil {
			return appErr
		}
		return nil
	})
	if appErr != nil {
		return DTO.JWTResponse{}, appErr
	}
	return jwtResponse, nil
}

func (s *Service) SignIn(ctx context.Context, request DTO.SignInRequest) (DTO.JWTResponse, *errs.AppError) {
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	adminUser, err := s.q.GetAdminUser(timeout, request.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DTO.JWTResponse{}, errs.Unauthorized(fmt.Errorf("no user with such email %w", err))
		}
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminUser.PasswordHash), []byte(request.Password))
	if err != nil {
		return DTO.JWTResponse{}, errs.Unauthorized(fmt.Errorf("wrong password %w", err))
	}
	jwtResponse, appErr := generateJWTResponse(s.c.JwtSecret, strconv.FormatInt(adminUser.ID, 10))
	if appErr != nil {
		return DTO.JWTResponse{}, appErr
	}
	return jwtResponse, nil
}

func (s *Service) TokenRefresh(ctx context.Context, request DTO.TokenRefreshRequest) (DTO.JWTResponse, *errs.AppError) {
	secret := s.c.JwtSecret
	withClaims, appErr := ParseToken(request.RefreshToken, secret)
	if appErr != nil {
		return DTO.JWTResponse{}, appErr
	}
	userID, err := withClaims.Claims.GetSubject()
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}

	jwtResponse, appErr := generateJWTResponse(secret, userID)
	if appErr != nil {
		return DTO.JWTResponse{}, appErr
	}
	return jwtResponse, nil
}

func newJWT(id string, exp, nbf time.Time) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%s-%d", id, time.Now().Nanosecond()),
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(exp),
		NotBefore: jwt.NewNumericDate(nbf),
	})
}

func generateJWTResponse(secret, id string) (DTO.JWTResponse, *errs.AppError) {
	accessToken := newJWT(id, time.Now().Add(accessExp), time.Now())
	signedAccessToken, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	refreshToken := newJWT(id, time.Now().Add(refreshExp), time.Now())
	signedRefreshToken, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	return mapJWTResponse(signedAccessToken, signedRefreshToken), nil
}

func ParseToken(token, secret string) (*jwt.Token, *errs.AppError) {
	withClaims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method == jwt.SigningMethodHS256 {
			return []byte(secret), nil
		}
		return nil, fmt.Errorf("wrong signing method - %s", token.Method.Alg())
	})
	if err != nil {
		return nil, errs.Unauthorized(err)
	}
	return withClaims, nil
}