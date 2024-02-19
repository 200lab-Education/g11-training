package usecase

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
