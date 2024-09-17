package utils

var (
	DATE_TIME_FORMAT = "2006-01-02 15:04:05"
	DATE_FORMAT      = "2006-01-02"

	// actor type for register and login
	User = "user"

	// verification type
	VerifyRegistration = "verify_registration"
	ForgotPassword     = "forgot_password"
	UpdateEmail        = "update_email"

	// reset password route
	ResetPassRoute   = "reset-password?token="
	VerifyEmailRoute = "verify-email?token="
	TypeRoute        = "&type="

	VerificationType = []string{VerifyRegistration, ForgotPassword, UpdateEmail}
)
