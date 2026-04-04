package auth

import (
	"fmt"

	"github.com/Aoladiy/go-with-tools/internal/errs"

	"github.com/golang-jwt/jwt/v5"
)

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
