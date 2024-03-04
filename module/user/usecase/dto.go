package usecase

import (
	"github.com/google/uuid"
	"my-app/common"
)

type EmailPasswordRegistrationDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type EmailPasswordLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponseDTO struct {
	AccessToken       string `json:"access_token"`
	AccessTokenExpIn  int    `json:"access_token_exp_in"`
	RefreshToken      string `json:"refresh_token"`
	RefreshTokenExpIn int    `json:"refresh_token_exp_in"`
}

type SingleImageDTO struct {
	Requester common.Requester `json:"-"`
	ImageId   uuid.UUID        `json:"image_id"`
}
