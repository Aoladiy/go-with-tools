package auth

import "go-with-tools/internal/DTO"

func mapJWTResponse(accessToken, refreshToken string) DTO.JWTResponse {
	return DTO.JWTResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}