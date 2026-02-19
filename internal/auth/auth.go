package auth

import (
	"context"
	"errors"
	"go-with-tools/internal/DTO"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/errs"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

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
	secret, isset := os.LookupEnv("JWT_SECRET")
	if !isset {
		return DTO.JWTResponse{}, errs.Internal(errors.New("cannot generate jwt token"))
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(adminUser.ID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		NotBefore: jwt.NewNumericDate(time.Now()),
	})
	signedAccessToken, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return DTO.JWTResponse{}, errs.Internal(err)
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
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
