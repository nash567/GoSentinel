package model

const (
	VerificationEmail = "verification email"
)

type Application struct {
	ID           string `json:"id"`
	Secret       string `json:"secret"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Status       string `json:"status"`
	IsVerified   bool   `json:"is_verified"`
	SecretViewed bool   `json:"secret_viewed"`
}

type MailData struct {
	URL      string `json:"url"`
	Template string `json:"template"`
}

type VerifyApplicationResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
type UpdateApplication struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	SecretViewed bool   `json:"secret_viewed"`
}
