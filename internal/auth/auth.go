package auth

import (
	"context"
	"errors"
	"fmt"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const UserId = "user_id"

type Service struct {
	q *queries.Queries
	p *pgxpool.Pool
}

func New(q *queries.Queries, p *pgxpool.Pool) *Service {
	return &Service{q: q, p: p}
}

func (s *Service) SignUp(ctx context.Context, request DTO.SignUpRequest) (DTO.JWTResponse, *errs.AppError) {
	if len(request.Password) < 8 {
		return DTO.JWTResponse{}, errs.BadRequest(errors.New("password must be at least 8 characters"))
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}

	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	tx, err := s.p.Begin(timeout)
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	defer tx.Rollback(timeout)
	adminUser, err := s.q.WithTx(tx).CreateAdminUser(timeout, queries.CreateAdminUserParams{
		Email:        request.Email,
		PasswordHash: string(password),
	})
	if err != nil {
		if pgErr, isUniqueViolation := errs.IsUniqueViolation(err); isUniqueViolation {
			return DTO.JWTResponse{}, errs.UniqueViolation(err, pgErr) // TODO handle case with unique email
		}
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	secret, isset := os.LookupEnv("JWT_SECRET") //TODO refactor to fill all env variables once on app startup
	if !isset {
		return DTO.JWTResponse{}, errs.Internal(errors.New("cannot generate jwt token"))
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%s-%d", strconv.FormatInt(adminUser.ID, 10), time.Now().Nanosecond()),
		Subject:   strconv.FormatInt(adminUser.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})
	signedAccessToken, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%s-%d", strconv.FormatInt(adminUser.ID, 10), time.Now().Nanosecond()),
		Subject:   strconv.FormatInt(adminUser.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	err = tx.Commit(timeout)
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	return DTO.JWTResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
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
	secret, isset := os.LookupEnv("JWT_SECRET") //TODO refactor to fill all env variables once on app startup
	if !isset {
		return DTO.JWTResponse{}, errs.Internal(errors.New("cannot generate jwt token"))
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%s-%d", strconv.FormatInt(adminUser.ID, 10), time.Now().Nanosecond()),
		Subject:   strconv.FormatInt(adminUser.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})
	signedAccessToken, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%s-%d", strconv.FormatInt(adminUser.ID, 10), time.Now().Nanosecond()),
		Subject:   strconv.FormatInt(adminUser.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	return DTO.JWTResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}

func (s *Service) TokenRefresh(ctx context.Context, request DTO.TokenRefreshRequest) (DTO.JWTResponse, *errs.AppError) {
	secret, isset := os.LookupEnv("JWT_SECRET") //TODO refactor to fill all env variables once on app startup
	if !isset {
		return DTO.JWTResponse{}, errs.Internal(errors.New("cannot generate jwt token"))
	}
	withClaims, err := jwt.ParseWithClaims(request.RefreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method == jwt.SigningMethodHS256 {
			return []byte(secret), nil
		}
		return nil, fmt.Errorf("wrong signing method - %s", token.Method.Alg())
	})
	if err != nil {
		return DTO.JWTResponse{}, errs.Unauthorized(err)
	}
	userID, err := withClaims.Claims.GetSubject()
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%s-%d", userID, time.Now().Nanosecond()),
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})
	signedAccessToken, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        fmt.Sprintf("%s-%d", userID, time.Now().Nanosecond()),
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	return DTO.JWTResponse{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil
}
