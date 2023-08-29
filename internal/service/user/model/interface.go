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
	RegisterUser(context.Context, User) error
	LoginUser(ctx context.Context, email string, password string) error
	GetUser(context.Context, Filter) (User, error)
	UpdateUser(context.Context, User) error
	DeleteUser(context.Context, string) error
}
