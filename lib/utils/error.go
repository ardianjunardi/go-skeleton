package utils

var (
	// General error
	EmptyData                          = "Data not found"
	ErrNotFoundPage                    = "Sorry. We couldn't find that page"
	ErrSystemError                     = "Something error with our system. Please contact our administrator"
	ErrInvalidTokenChannel             = "Invalid token channel"
	ErrConfigKeyNotFound               = "Config ['app.key'] doesn't exists"
	ErrEmailAlreadyRegistered          = "Your email has been registered. Please change to a new email"
	ErrEmailNotVerified                = "Email not verified"
	ErrInvalidEmailPassword            = "Email/password is incorrect"
	ErrGeneratingJWT                   = "Error generating JWT token"
	ErrGettingVerificationsData        = "Error getting data from verifications table"
	ErrBeginningTransaction            = "Error beginning transaction"
	ErrInvalidToken                    = "Token is invalid"
	ErrTokenExpired                    = "Token is expired"
	ErrTokenUsed                       = "Token has been used"
	ErrMarkingToken                    = "Error marking token as used"
	ErrCommittingTransaction           = "Error committing transaction"
	ErrSendingResetPasswordEmail       = "Error sending email for reset password"
	ErrAddingResetPasswordVerification = "Error adding verification for reset password"
	ErrSendingVerifyEmail              = "Error sending email for verify email"
	ErrSendingForgotPasswordEmail      = "Error sending email for forgot password"
	ErrSendingUpdateEmail              = "Error sending email for update email"
	ErrInvalidSendingEmailType         = "Type must be one of the following: (verify_registration | forgot_password | update_email)"
	ErrInvalidTypeQueryParameter       = "Type query parameter is missing"
	ErrPasswordMismatch                = "Password does not match"
	ErrHashingPassword                 = "Error hashing the new password"
	ErrInvalidTypeError                = "Incorrect error type provided for 'err' parameter. It must be an instance of 'validator.ValidationErrors' or 'error'"

	// Error for module AMQP
	ErrConnectAMQP           = "Can't connect to AMQP"
	ErrCreateChannelAMQP     = "Can't create a amqpChannel"
	ErrContentTypeNotAllowed = "Content type is not allowed"

	// Error for module user
	ErrGettingUserData                = "Error getting user data"
	ErrInsertingUser                  = "Error inserting user"
	ErrUpdatingUserPassword           = "Error updating user's password"
	ErrFetchingUserPassword           = "Error fetching user's current password"
	ErrUpdatingUserEmail              = "Error updating user email"
	ErrUpdatingUserEmailStatus        = "Error updating user email status"
	ErrUpdatingUserProfile            = "Error updating user profile"
	ErrRetrievingUserByUserIdentifier = "Error retrieving user by user identifier"
	ErrGettingUserByEmail             = "Error getting user by email"

	// Error for module user address
	ErrGettingUserAddresses      = "Error getting user addresses by user ID"
	ErrScanningUserAddresses     = "Error scanning user addresses"
	ErrIteratingUserAddresses    = "Error iterating over user addresses"
	ErrInsertingUserAddress      = "Error inserting user address"
	ErrCheckingAddressIdentifier = "Error checking address identifier"
	ErrUpdatingUserAddress       = "Error updating user address"
	ErrDeletingUserAddress       = "Error deleting user address"

	// Error for module setting
	ErrCountingListSetting  = "Error counting list setting"
	ErrGettingListSetting   = "Error getting list setting"
	ErrScanningListSetting  = "Error scanning list setting"
	ErrAddingSetting        = "Error adding setting"
	ErrGettingSettingByCode = "Error getting setting by code"
	ErrUpdatingSetting      = "Error updating setting"
	ErrGettingSettingByKey  = "Error getting setting by key"
)
