package request

type UpdateProfileUserReq struct {
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name"`
	AvatarUrl   string `json:"avatar_url"`
	Description string `json:"description"`
}

type UpdatePasswordReq struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
