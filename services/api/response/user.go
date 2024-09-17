package response

type UserRes struct {
	UserIdentifier string `json:"user_identifier"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	AvatarURL      string `json:"avatar_url"`
	Password       string `json:"password"`
	IsVerified     bool   `json:"is_verify"`
	CreatedDate    string `json:"created_date"`
	UpdatedDate    string `json:"updated_date"`
	DeletedDate    string `json:"deleted_date"`
}

type UserProfileRes struct {
	UserIdentifier string `json:"user_identifier"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	AvatarURL      string `json:"avatar_url"`
	IsVerified     bool   `json:"is_verify"`
	CreatedDate    string `json:"created_date"`
	UpdatedDate    string `json:"updated_date"`
}

type UserAddressRes struct {
	AddressIdentifier string `json:"address_identifier"`
	Title             string `json:"title"`
	FullAddress       string `json:"full_address"`
	CreatedDate       string `json:"created_date"`
	UpdatedDate       string `json:"updated_date"`
	DeletedDate       string `json:"deleted_date"`
}
