package model

const (
	VerificationEmail = "verification email"
)

type Application struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Status       string `json:"status"`
	Password     string `json:"password"`
	IsVerified   bool   `json:"is_verified"`
	SecretViewed bool   `json:"secret_viewed"`
}

type ApplicationSecret struct {
	ApplicationID     string `json:"id"`
	ApplicationSecret string `json:"secret"`
}

type MailData struct {
	URL      string `json:"url"`
	Template string `json:"template"`
}

type VerifyApplicationResponse struct {
	ApplicationID     string `json:"client_id"`
	ApplicationSecret string `json:"client_secret"`
}
type UpdateApplication struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
