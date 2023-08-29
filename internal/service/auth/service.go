package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/nash567/GoSentinel/internal/service/auth/config"
	"github.com/nash567/GoSentinel/internal/service/auth/model"
)

type Service struct {
	config *config.Config
	repo   model.Repository
}

func NewService(config *config.Config, repo model.Repository) *Service {
	return &Service{
		config: config,
		repo:   repo,
	}

}
func (s *Service) CreateApplicationIdentity(ctx context.Context, applicationId string) error {
	credentials, err := s.GenerateCredentials(s.config.SecretLength)
	if err != nil {
		return fmt.Errorf("generating credentials:%w", err)
	}
	encryptedSecret, err := s.EncryptData(credentials.ApplicationSecret, s.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf(" encrypting secret:%w", err)
	}
	err = s.repo.CreateApplicationIdentity(ctx, model.ApplicationIdentity{
		ID:            credentials.ApplicationID,
		ApplicationID: applicationId,
		Secret:        encryptedSecret,
	})
	if err != nil {
		return fmt.Errorf("creating application identity:%w", err)
	}

	return nil

}

func (s *Service) GetApplicationIdentity(ctx context.Context, applicationID string) (*model.ApplicationIdentity, error) {
	identity, err := s.repo.GetApplicationIdentity(ctx, applicationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("applicationID not valid: %w", err)
		}
		return nil, fmt.Errorf("get application identity: %w", err)
	}
	if identity.SecretViewed {
		return nil, fmt.Errorf("secret already viewed: %w", err)
	}

	secretBytes, err := s.DecryptData(identity.Secret, s.config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}
	identity.Secret = string(secretBytes)
	return identity, nil
}

func (s *Service) UpdateApplicationIdentity(ctx context.Context, applicationID string) error {
	err := s.repo.UpdateApplicationIdentity(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("failed to updated application identity: %w", err)
	}

	return nil
}

func (s *Service) VerifyApplicationIdentity(ctx context.Context, credentials model.Credentials) (bool, error) {
	identity, err := s.repo.GetApplicationIdentity(ctx, credentials.ApplicationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("applicationID not valid: %w", err)
		}
		return false, fmt.Errorf("get application identity: %w", err)
	}

	secretBytes, err := s.DecryptData(identity.Secret, s.config.EncryptionKey)
	if err != nil {
		return false, fmt.Errorf("failed to decrypt secret: %w", err)
	}
	fmt.Println(string(secretBytes))
	fmt.Println("application secret", credentials.ApplicationSecret)
	if strings.Compare(string(secretBytes), credentials.ApplicationSecret) != 0 {
		return false, fmt.Errorf("wrong secret key")
	} else {
		return true, nil
	}

}
