package models

import "github.com/golang-jwt/jwt"

type JWTClaim struct {
	UserID uint   `json:"id"`
	Email  string `json:"email"`
	Role string `json:"role"`
	jwt.StandardClaims
}