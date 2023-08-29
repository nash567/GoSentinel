package user

import (
	"context"
	"fmt"

	"github.com/nash567/GoSentinel/internal/service/user/model"
)

type Service struct {
	repo model.Repository
}

func NewService(repo model.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) RegisterUser(ctx context.Context, user model.User, applicationID string) error {
	err := s.repo.RegisterUser(ctx, user, applicationID)
	if err != nil {
		return fmt.Errorf("register user: %w", err)
	}
	return nil
}
func (s *Service) GetUser(ctx context.Context, filter model.Filter, applicationID string) (*model.User, error) {
	user, err := s.repo.GetUser(ctx, filter, applicationID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}
func (s *Service) UpdateUser(ctx context.Context, user model.UpdateUser, applicationID string) error {
	err := s.repo.UpdateUser(ctx, user, applicationID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}
func (s *Service) DeleteUser(ctx context.Context, userID string, applicationID string) error {
	err := s.repo.DeleteUser(ctx, userID, applicationID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
