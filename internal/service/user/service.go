package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	authConfig "github.com/nash567/GoSentinel/internal/service/auth/config"
	authModel "github.com/nash567/GoSentinel/internal/service/auth/model"
	"github.com/nash567/GoSentinel/internal/service/user/model"
)

type Service struct {
	repo       model.Repository
	authSvc    authModel.Service
	authConfig authConfig.Config
}

func NewService(repo model.Repository, authSvc authModel.Service, authConfig authConfig.Config) *Service {

	return &Service{
		repo:       repo,
		authSvc:    authSvc,
		authConfig: authConfig,
	}

}

func (s *Service) RegisterUser(ctx context.Context, user model.User) error {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return fmt.Errorf("claims not found in context")
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generating uuid: %w", err)
	}
	user.ID = id.String()

	password, err := s.authSvc.EncryptData(user.Password, s.authConfig.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting password: %w", err)
	}
	user.Password = password

	err = s.repo.RegisterUser(ctx, user, aws.StringValue(claims.ApplicationID))
	if err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	return nil
}
func (s *Service) GetUser(ctx context.Context, filter model.Filter) (*model.User, error) {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	user, err := s.repo.GetUser(ctx, filter, aws.StringValue(claims.ApplicationID))
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}
func (s *Service) UpdateUser(ctx context.Context, user model.UpdateUser) error {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return fmt.Errorf("claims not found in context")
	}
	err := s.repo.UpdateUser(ctx, user, aws.StringValue(claims.ApplicationID))
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}
func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return fmt.Errorf("claims not found in context")
	}
	err := s.repo.DeleteUser(ctx, userID, aws.StringValue(claims.ApplicationID))
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (*string, error) {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	user, err := s.repo.GetUser(ctx, model.Filter{
		Email: []string{email},
	}, *claims.ApplicationJWTClaims.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("get user : %w", err)
	}

	decryptedPassword, err := s.authSvc.DecryptData(user.Password, s.authConfig.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt data : %w", err)
	}
	if strings.Compare(string(decryptedPassword), password) != 0 {
		return nil, fmt.Errorf("user authentication failed")
	}

	token, err := s.authSvc.GenerateJWtToken(authModel.Claims{
		UserJWTClaims: &authModel.UserJWTClaims{
			UserEmail: aws.String(user.Email),
			UserID:    aws.String(user.ID),
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * s.authConfig.UserJWTExpiry)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("generate jwt token: %w", err)
	}
	return &token, nil
}
