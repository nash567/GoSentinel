package model

import (
	"context"

	"google.golang.org/grpc"
)

type Service interface {
	CreateApplicationIdentity(ctx context.Context, applicationId string) error
	GetApplicationIdentity(ctx context.Context, applicationID string) (*ApplicationIdentity, error)
	UpdateApplicationIdentity(ctx context.Context, applicationID string) error
	VerifyApplicationIdentity(ctx context.Context, credentials Credentials) (bool, error)
	GenerateJWtToken(claims Claims) (string, error)
	VerifyJWTToken(token string) (*Claims, error)
	GenerateSecret(length int) (string, error)
	GetApplicationToken(ctx context.Context, credentials Credentials) (*string, error)
	AuthenticationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	EncryptData(data string, key string) (string, error)
	DecryptData(encryptedData string, key string) ([]byte, error)
}

type Repository interface {
	CreateApplicationIdentity(ctx context.Context, identity ApplicationIdentity) error
	GetApplicationIdentity(ctx context.Context, applicationID string) (*ApplicationIdentity, error)
	UpdateApplicationIdentity(ctx context.Context, applicationID string) error
}
