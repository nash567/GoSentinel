package model

import (
	"context"
)

type Service interface {
	RegisterUser(context.Context, User) error
	LoginUser(ctx context.Context, email, password string) (*string, error)
	GetUser(context.Context, Filter) (*User, error)
	UpdateUser(context.Context, UpdateUser) error
	DeleteUser(context.Context, string) error
}

type Repository interface {
	GetUser(context.Context, Filter, string) (*User, error)
	RegisterUser(context.Context, User, string) error
	UpdateUser(context.Context, UpdateUser, string) error
	DeleteUser(context.Context, string, string) error
}
