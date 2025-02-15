package request

type LoginUserReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserReq struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type RequestVerifyEmailReq struct {
	Email string `json:"email" validate:"required"`
	Type  string `json:"type" validate:"required"`
}

type VerifyEmailReq struct {
	Email string `json:"email" validate:"required"`
	Token string `json:"token" validate:"required"`
}

type ResetPasswordReq struct {
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
