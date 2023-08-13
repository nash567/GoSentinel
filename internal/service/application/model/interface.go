package model

import (
	"context"
)

type Service interface {
	SendVerifcationNotification(ctx context.Context, email string, name string) (*string, error)
	VerifyApplication(ctx context.Context, key string) error
	GetApplicationSecret(ctx context.Context, token string) (*Application, error)
}

type Repository interface {
	RegisterApplication(ctx context.Context, application Application) error
	GetApplication(ctx context.Context, filter *Filter) (*Application, error)
	UpdateApplication(ctx context.Context, application *UpdateApplication) error
}
