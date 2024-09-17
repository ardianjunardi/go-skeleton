package response

type LoginUserRes struct {
	Token          string `json:"token"`
	UserIdentifier string `json:"user_identifier"`
	AvatarURL      string `json:"avatar_url"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	ExpiredAt      string `json:"expired_at"`
	ActorType      string `json:"actor_type"`
	CreatedDate    string `json:"created_date"`
}

type RegisterUserRes struct {
	UserIdentifier string `json:"user_identifier"`
	Email          string `json:"email"`
	CreatedDate    string `json:"created_date"`
}
