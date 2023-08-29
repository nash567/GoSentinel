package model

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	ApplicationID string `json:"application_id"`
	jwt.RegisteredClaims
}

type Credentials struct {
	ApplicationID     string `json:"id"`
	ApplicationSecret string `json:"secret"`
}

type ApplicationIdentity struct {
	ID            string `json:"id"`
	ApplicationID string `json:"application_id"`
	Secret        string `json:"secret"`
	SecretViewed  bool   `json:"secret_viewed"`
}
type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ContextKeyJWTClaims = contextKey("jwtClaims")
)

func GetJWTClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(ContextKeyJWTClaims).(*Claims)
	return claims, ok
}
