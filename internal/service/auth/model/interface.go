package model

type Service interface {
	GenerateJWtToken(claims Claims) (string, error)
	VerifyJWTToken(token string) (*Claims, error)
	GenerateCredentials(length int) (*Credentials, error)
	EncryptData(data string, key string) (string, error)
	DecryptData(encryptedData string, key string) ([]byte, error)
	GenerateSecret(length int) (string, error)
}
