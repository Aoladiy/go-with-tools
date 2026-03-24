package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Aoladiy/go-with-tools/internal/DTO"
	"github.com/Aoladiy/go-with-tools/internal/config"
	"github.com/Aoladiy/go-with-tools/internal/database/queries"
	"github.com/Aoladiy/go-with-tools/internal/errs"
	"github.com/Aoladiy/go-with-tools/internal/helpers"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessExp  = time.Minute * 15
	refreshExp = time.Hour * 24 * 7
	SignedOut  = "signed-out-token-"
)

type Service struct {
	q   *queries.Queries
	rdb *redis.Client
	p   *pgxpool.Pool
	c   config.Config
}

func New(q *queries.Queries, rdb *redis.Client, p *pgxpool.Pool, c config.Config) *Service {
	return &Service{q: q, rdb: rdb, p: p, c: c}
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
	withClaims, appErr := ParseToken(request.RefreshToken, s.c.JwtSecret)
	if appErr != nil {
		return DTO.JWTResponse{}, appErr
	}
	isTokenSignedOut, err := IsTokenSignedOut(ctx, s.rdb, withClaims.Raw)
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	if isTokenSignedOut {
		return DTO.JWTResponse{}, errs.Unauthorized(fmt.Errorf("token is in signed out tokens cache"))
	}
	userID, err := withClaims.Claims.GetSubject()
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}

	jwtResponse, appErr := generateJWTResponse(s.c.JwtSecret, userID)
	if appErr != nil {
		return DTO.JWTResponse{}, appErr
	}
	return jwtResponse, nil
}

func (s *Service) SignOut(ctx context.Context, request DTO.SignOutRequest) *errs.AppError {
	accessToken, appErr := ParseToken(request.AccessToken, s.c.JwtSecret)
	if appErr != nil {
		return appErr
	}
	refreshToken, appErr := ParseToken(request.RefreshToken, s.c.JwtSecret)
	if appErr != nil {
		return appErr
	}
	accessTokenExp, err := accessToken.Claims.GetExpirationTime()
	if err != nil {
		return errs.BadRequest(err)
	}
	refreshTokenExp, err := refreshToken.Claims.GetExpirationTime()
	if err != nil {
		return errs.BadRequest(err)
	}
	set := s.rdb.Set(ctx, SignedOut+accessToken.Raw, true, accessTokenExp.Time.Sub(time.Now()))
	if set.Err() != nil {
		return errs.Internal(set.Err())
	}
	set = s.rdb.Set(ctx, SignedOut+refreshToken.Raw, true, refreshTokenExp.Time.Sub(time.Now()))
	if set.Err() != nil {
		return errs.Internal(set.Err())
	}
	return nil
}
