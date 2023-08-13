package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
	notification "github.com/nash567/GoSentinel/internal/notifications/email/model"
	"github.com/nash567/GoSentinel/internal/service/application/config"
	authConfig "github.com/nash567/GoSentinel/internal/service/auth/config"

	"github.com/nash567/GoSentinel/internal/service/application/model"
	authModel "github.com/nash567/GoSentinel/internal/service/auth/model"
	"github.com/nash567/GoSentinel/pkg/cache"
)

const (
	baseURL = `http://localhost:8080/v1/verify`
)

type Service struct {
	repo       model.Repository
	config     *config.Config
	authConfig authConfig.Config
	mailSvc    notification.Service
	cacheSvc   cache.Cache
	authSvc    authModel.Service
}

func NewService(
	config *config.Config,
	mailSvc notification.Service,
	repo model.Repository,
	cacheSvc cache.Cache,
	authConfig authConfig.Config,
	authSvc authModel.Service,
) *Service {
	return &Service{
		config:     config,
		mailSvc:    mailSvc,
		repo:       repo,
		cacheSvc:   cacheSvc,
		authConfig: authConfig,
		authSvc:    authSvc,
	}
}

func (s *Service) SendVerifcationNotification(ctx context.Context, email, name string) (*string, error) {
	// check if email already exists
	application, err := s.repo.GetApplication(ctx, &model.Filter{
		Email: []string{email},
	})
	if err == nil && application != nil {
		return nil, fmt.Errorf("user exist already")
	}

	//generate secret
	key, err := s.authSvc.GenerateSecret(s.authConfig.SecretLength)
	if err != nil {
		return nil, fmt.Errorf("svc:application -> generateToken: %w", err)
	}

	// set secret in cache
	err = s.cacheSvc.Set(ctx, cache.NewKeyValWithExpiry(key, model.Application{
		Email: email,
		Name:  name,
	}, s.config.VerificationExpiry*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("svc:application -> cacheSvc.Set: %w", err)
	}

	//get template
	template, err := s.getTemplate(model.MailData{
		URL:      fmt.Sprintf("%s/%s", baseURL, key),
		Template: s.config.VerificationTemplate,
	})
	if err != nil {
		return nil, fmt.Errorf("svc:application -> getTemplate: %w", err)
	}

	// send email
	err = s.mailSvc.Send(ctx, notification.NewMail([]string{email}, model.VerificationEmail, template))
	if err != nil {
		return nil, fmt.Errorf("svc:application -> send: %w", err)
	}

	//generate token
	token, err := s.authSvc.GenerateJWtToken(authModel.Claims{
		Email: email,
		Name:  name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * s.authConfig.JWTExpiration)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("svc:application -> generateToken: %w", err)
	}

	return &token, nil
}

func (s *Service) VerifyApplication(ctx context.Context, key string) error {
	// get value
	value, err := s.cacheSvc.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("svc:application -> cacheSvc get: %w", err)
	}

	// unmarshal value
	var application model.Application
	err = json.Unmarshal([]byte(value), &application)
	if err != nil {
		return fmt.Errorf("unmarshal zip status from cache: %w", err)
	}

	creds, err := s.authSvc.GenerateCredentials(s.authConfig.SecretLength)
	if err != nil {
		return fmt.Errorf("svc:application -> generateCredentials: %w", err)
	}
	secret, err := s.authSvc.EncryptData(creds.Secret, s.authConfig.EncryptionKey)
	if err != nil {
		return fmt.Errorf("svc:application -> encrypt: %w", err)
	}

	application.ID = creds.ID
	application.Secret = secret
	application.Status = "active"
	application.IsVerified = true
	err = s.repo.RegisterApplication(ctx, application)
	if err != nil {
		return fmt.Errorf("svc:application -> RegisterApplication: %w", err)
	}
	return nil
}

func (s *Service) getTemplate(mailData model.MailData) (string, error) {
	t := template.New(filepath.Base(mailData.Template))
	t, err := t.ParseFiles(mailData.Template)
	if err != nil {
		return "", fmt.Errorf("parse template files: %w", err)
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, mailData)
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	result := tpl.String()
	return result, nil
}

// func (s *Service) generateURL(baseURL string, token string) (string, error) {
// 	u, err := url.Parse(baseURL)
// 	if err != nil {

// 		return "", fmt.Errorf("svc:application -> url.Parse: %w", err)
// 	}
// 	q := u.Query()
// 	q.Set("token", token) // Attach the token as a query parameter
// 	u.RawQuery = q.Encode()
// 	return u.String(), nil
// }

func (s *Service) AuthenticateApplication(ctx context.Context, input model.Application) error {
	application, err := s.repo.GetApplication(ctx, &model.Filter{
		Email: []string{input.Email},
	})
	if err == nil && application != nil {
		return fmt.Errorf("user exist already")
	}

	secret, err := s.authSvc.DecryptData(application.Secret, s.authConfig.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error decrypting :%w", err)
	}

	if string(secret) == input.Secret {
		return nil
	}

	return fmt.Errorf("user unauthenticated")
}

func (s *Service) GetApplicationSecret(ctx context.Context, token string) (*model.Application, error) {
	claims, err := s.authSvc.VerifyJWTToken(token)
	if err != nil {
		return nil, fmt.Errorf("svc:application -> verifyToken: %w", err)
	}

	application, err := s.repo.GetApplication(ctx, &model.Filter{
		Email: []string{claims.Email},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting application :%w", err)
	}

	if application.SecretViewed {
		return nil, fmt.Errorf("secret already viewed")
	} else {
		err := s.repo.UpdateApplication(ctx, &model.UpdateApplication{
			ID: application.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("error updating  application :%w", err)
		}
		return application, nil

	}
}
