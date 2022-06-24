package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserID int `json:"userID"`
	jwt.StandardClaims
}
