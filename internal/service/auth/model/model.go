package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}

type Credentials struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}
