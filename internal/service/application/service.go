package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
		return nil, fmt.Errorf("application exist already")
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
		ApplicationJWTClaims: &authModel.ApplicationJWTClaims{
			ApplicationEmail: aws.String(email),
			Name:             aws.String(name),
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * s.authConfig.VerificationJWTExpiration)),
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
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generating uuid: %v", err)
	}
	application.ID = id.String()
	application.Status = "active"
	application.IsVerified = true
	err = s.repo.RegisterApplication(ctx, application)
	if err != nil {
		return fmt.Errorf("svc:application -> RegisterApplication: %w", err)
	}
	return nil
}
func (s *Service) CreateApplicationPassword(ctx context.Context, application *model.UpdateApplication) error {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return fmt.Errorf("claims not found in context")
	}
	app, err := s.repo.GetApplication(ctx, &model.Filter{
		Email: []string{aws.StringValue(claims.ApplicationJWTClaims.ApplicationEmail)},
	})
	if err != nil {
		return fmt.Errorf("error getting application :%w", err)
	}
	encryptedPassword, err := s.authSvc.EncryptData(application.Password, s.authConfig.EncryptionKey)
	if err != nil {
		return fmt.Errorf("error encrypting password :%w", err)
	}
	application.Password = encryptedPassword
	application.ID = app.ID
	err = s.repo.UpdateApplication(ctx, application)
	if err != nil {
		return fmt.Errorf("failed to update application password:%w", err)
	}

	return nil
}

func (s *Service) LoginApplication(ctx context.Context, email, password string) (*string, error) {
	application, err := s.repo.GetApplication(ctx, &model.Filter{
		Email: []string{email},
	})
	if err != nil {
		return nil, fmt.Errorf("get application : %w", err)
	}

	decryptedPassword, err := s.authSvc.DecryptData(aws.StringValue(application.Password), s.authConfig.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt data : %w", err)
	}
	if strings.Compare(string(decryptedPassword), password) != 0 {
		return nil, fmt.Errorf("application authentication failed")
	}

	token, err := s.authSvc.GenerateJWtToken(authModel.Claims{
		ApplicationJWTClaims: &authModel.ApplicationJWTClaims{
			ApplicationEmail: aws.String(application.Email),
			ApplicationID:    aws.String(application.ID),
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * s.authConfig.ApplicationJWTExpiry)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("generate jwt token: %w", err)
	}
	return &token, nil
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

func (s *Service) CreateApplicationIdentity(ctx context.Context) error {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return fmt.Errorf("claims not found in context")
	}
	application, err := s.repo.GetApplication(ctx, &model.Filter{
		Email: []string{aws.StringValue(claims.ApplicationJWTClaims.ApplicationEmail)},
	})
	if err != nil {
		return fmt.Errorf("error getting application :%w", err)
	}
	err = s.authSvc.CreateApplicationIdentity(ctx, application.ID)
	if err != nil {
		return fmt.Errorf("error creating application identity: %w", err)
	}
	return nil

}

func (s *Service) GetApplicationIdentity(ctx context.Context) (*model.ApplicationSecret, error) {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	application, err := s.repo.GetApplication(ctx, &model.Filter{
		Email: []string{aws.StringValue(claims.ApplicationJWTClaims.ApplicationEmail)},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting application :%w", err)
	}

	applicationIdentity, err := s.authSvc.GetApplicationIdentity(ctx, application.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting application identity :%w", err)
	}
	err = s.authSvc.UpdateApplicationIdentity(ctx, application.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating  application identity :%w", err)
	}

	return &model.ApplicationSecret{
		ApplicationID:     applicationIdentity.ApplicationID,
		ApplicationSecret: applicationIdentity.Secret,
	}, nil

}
