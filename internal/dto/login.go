package dto

type LoginRequest struct {
	DocumentNumber string `json:"documentNumber" validate:"required"`
	Password       string `json:"password" validate:"required"`
}

type Login struct {
	AccessToken string `json:"access_token,omitempty"`
}
