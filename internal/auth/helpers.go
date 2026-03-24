package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Aoladiy/go-with-tools/internal/DTO"
	"github.com/Aoladiy/go-with-tools/internal/errs"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

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

func IsTokenSignedOut(ctx context.Context, rdb *redis.Client, token string) (bool, error) {
	err := rdb.Get(ctx, SignedOut+token).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	} else if errors.Is(err, redis.Nil) {
		return false, nil
	}
	return true, nil
}
