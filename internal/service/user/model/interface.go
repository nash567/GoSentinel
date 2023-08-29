package model

import (
	"context"
)

type Service interface {
	common
}

type Repository interface {
	common
}

type common interface {
	RegisterUser(context.Context, User, string) error
	LoginUser(ctx context.Context, email string, password string, applicationID string) error
	GetUser(context.Context, Filter, string) (*User, error)
	UpdateUser(context.Context, UpdateUser, string) error
	DeleteUser(context.Context, string, string) error
}
