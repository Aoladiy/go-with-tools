package auth

import (
	"github.com/Aoladiy/go-with-tools/gen"
)

func mapJWTResponse(accessToken, refreshToken string) gen.JWTResponse {
	return gen.JWTResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
