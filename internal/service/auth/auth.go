package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nash567/GoSentinel/internal/service/auth/model"
)

func (s *Service) GenerateJWtToken(claims model.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.config.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("svc:application -> generateToken->SignedString :%w", err)
	}
	return tokenString, nil
}
func (s *Service) VerifyJWTToken(token string) (*model.Claims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Token parsing error: %w", err)
	}
	if parsedToken == nil {
		return nil, fmt.Errorf("invalid Token")

	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Token error")
	}

	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return nil, errors.New("Token expired")
	}

	return &model.Claims{
		Email: claims["email"].(string),
	}, nil

}

func (s *Service) GenerateCredentials(length int) (*model.Credentials, error) {
	creds := &model.Credentials{}
	secret, err := s.GenerateSecret(s.config.SecretLength)
	if err != nil {
		return nil, fmt.Errorf("generateSecret: %v", err)
	}
	id, err := s.generateUID()
	if err != nil {
		return nil, fmt.Errorf("generating uuid: %v", err)
	}
	creds.ApplicationID = id
	creds.ApplicationSecret = secret

	return creds, nil
}
func (s *Service) generateUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generating uuid: %v", err)
	}

	return id.String(), nil
}
func (s *Service) GenerateSecret(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "nil", fmt.Errorf("reading bytes: %v", err)
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
func (s *Service) EncryptData(data string, key string) (string, error) {
	AESkey, err := generateAESKey(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(AESkey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len([]byte(data)))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
func generateAESKey(passphrase string) ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(passphrase))
	return hash.Sum(nil), nil
}

func (s *Service) DecryptData(encryptedData string, key string) ([]byte, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	AESkey, err := generateAESKey(key)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(AESkey)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func (s *Service) GetApplicationToken(ctx context.Context, credentials model.Credentials) (*string, error) {
	verified, err := s.VerifyApplicationIdentity(ctx, credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to verify application identity :%w", err)
	}

	if verified {
		token, err := s.GenerateJWtToken(model.Claims{
			ApplicationID: credentials.ApplicationID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.ApplicationJWTExpiry * time.Minute)),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate token :%w", err)
		}

		return &token, nil

	}

	return nil, fmt.Errorf("application not valid")
}
