package model

import (
	"context"
	"errors"
	"go-skeleton/bootstrap"
	"go-skeleton/lib/mail"
	"go-skeleton/lib/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func (c *Contract) GenerateTokenJWT(userIdentifier, actorType, email string) (string, int64, error) {
	var (
		token string
		expAt int64
		err   error
	)

	key := c.Config.GetString("app.key")
	if len(key) == 0 {
		return token, expAt, errors.New(utils.ErrConfigKeyNotFound)
	}
	// 7 month duration
	expAt = time.Now().UTC().AddDate(7, 0, 0).Unix()
	claims := &bootstrap.CustomUserClaims{
		UserIdentifier: userIdentifier,
		Email:          email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    actorType,
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = rawToken.SignedString([]byte(key))
	if err != nil {
		return token, expAt, err
	}
	return token, expAt, nil
}

func (c *Contract) RegisterUser(db *pgxpool.Pool, ctx context.Context, firstName, lastName, email, password string) (string, string, error) {
	var (
		err           error
		id            int64
		userInsertSQL string

		// Generate User Identifier
		userIdentifier = utils.GeneratePrefixCode(utils.UserPrefix)

		// Replace hashPassword with your actual password hashing function
		passwordHash, _ = bcrypt.GenerateFromPassword([]byte(password), 14)
	)

	// Insert user data into 'users' table
	userInsertSQL = `INSERT INTO users (user_identifier, first_name, last_name, email, password, is_verify, created_date) 
        VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = db.QueryRow(ctx, userInsertSQL, userIdentifier, firstName, lastName, email, passwordHash, false, time.Now().In(time.UTC)).Scan(&id)
	if err != nil {
		// Handle specific error cases
		switch {
		case strings.Contains(err.Error(), "users_email_key"):
			return userIdentifier, email, errors.New(utils.ErrEmailAlreadyRegistered)
		// Add other specific error cases here if needed
		default:
			return userIdentifier, email, c.errHandler("model.RegisterUser", err, utils.ErrInsertingUser)
		}
	}

	return userIdentifier, email, nil
}

func (c *Contract) UserLogin(db *pgxpool.Pool, ctx context.Context, email, password string) (UserEnt, int64, string, error) {
	var (
		err      error
		userData UserEnt
		expAt    int64
		jwtToken string
	)

	// Get user data by email
	userData, err = c.GetUserByEmail(db, ctx, email)
	if err != nil {
		if err.Error() == utils.EmptyData {
			return userData, expAt, jwtToken, errors.New(utils.ErrInvalidEmailPassword)
		}
		return userData, expAt, jwtToken, c.errHandler("model.UserLogin", err, err.Error())
	}

	// If email is not verified
	if !userData.IsVerified {
		return userData, expAt, jwtToken, errors.New(utils.ErrEmailNotVerified)
	}

	dataPassword := []byte(userData.Password)
	err = bcrypt.CompareHashAndPassword(dataPassword, []byte(password))
	if err != nil {
		return userData, expAt, jwtToken, errors.New(utils.ErrInvalidEmailPassword)
	}

	// Generate JWT Token
	jwtToken, expAt, err = c.GenerateTokenJWT(userData.UserIdentifier, utils.User, email)
	if err != nil {
		return userData, expAt, jwtToken, c.errHandler("model.UserLogin", err, utils.ErrGeneratingJWT)
	}

	return userData, expAt, jwtToken, nil
}

func (c *Contract) CheckTokenAndExpiration(db *pgxpool.Pool, ctx context.Context, verificationType, actorType, token string) (UserEnt, int64, string, error) {
	var (
		jwtToken                string
		expAt                   int64
		checkVerificationSQL    string
		updateVerificationSQL   string
		email                   string
		expiredDateVerification time.Time
		isUsedData              bool
		dataUser                UserEnt
		err                     error
	)

	// Check if token exists and is not used
	checkVerificationSQL = `SELECT email, expired_date, is_used FROM verifications WHERE token = $1 AND verification_type = $2`
	err = db.QueryRow(ctx, checkVerificationSQL, token, verificationType).Scan(&email, &expiredDateVerification, &isUsedData)
	if err != nil {
		return dataUser, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrGettingVerificationsData)
	}

	// Check if token is used
	if isUsedData {
		return dataUser, expAt, jwtToken, errors.New(utils.ErrTokenUsed)
	}

	// Check if token is expired
	if expiredDateVerification.Before(time.Now().UTC()) {
		return dataUser, expAt, jwtToken, errors.New(utils.ErrTokenExpired)
	}

	if actorType == utils.User {
		// Get data user for create jwt token
		dataUser, err = c.GetUserByEmail(db, ctx, email)
		if err != nil {
			return dataUser, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiration", err, utils.ErrGettingUserData)
		}
	}
	// Generate JWT Token
	jwtToken, expAt, err = c.GenerateTokenJWT(dataUser.UserIdentifier, utils.User, email)
	if err != nil {
		return dataUser, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiration", err, utils.ErrGeneratingJWT)
	}
	// Start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return dataUser, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiration", err, utils.ErrBeginningTransaction)
	}

	// Update the verification record to mark it as used
	updateVerificationSQL = "UPDATE verifications SET is_used = $1 WHERE token = $2"
	_, err = tx.Exec(ctx, updateVerificationSQL, true, token)
	if err != nil {
		tx.Rollback(ctx)
		return dataUser, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrMarkingToken)
	}

	switch verificationType {
	// Check is type of update email
	case utils.UpdateEmail:
		sql := `
			UPDATE users
			SET email = $1 , is_verify = $2
			WHERE user_identifier = $3
		`

		_, err := tx.Exec(ctx, sql, email, true, dataUser.UserIdentifier)
		if err != nil {
			tx.Rollback(ctx)
			return dataUser, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrUpdatingUserEmail)
		}
	// Check is type of verify email
	case utils.VerifyRegistration:
		sql := `
			UPDATE users
			SET is_verify = $1
			WHERE user_identifier = $2
		`

		_, err := tx.Exec(ctx, sql, true, dataUser.UserIdentifier)
		if err != nil {
			tx.Rollback(ctx)
			return dataUser, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrUpdatingUserEmailStatus)
		}
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return dataUser, expAt, jwtToken, c.errHandler("model.CheckTokenAndExpiredForgotPassword", err, utils.ErrCommittingTransaction)
	}

	return dataUser, expAt, jwtToken, nil
}

func (c *Contract) RequestForgotPassword(db *pgxpool.Pool, ctx context.Context, email string) error {
	var (
		err error
		// Generate Token 50 digits
		token, _ = utils.Generate(`[a-zA-Z0-9]{50}`)

		// Forgot password route
		linkNewPass = c.Config.GetString("web_url") + utils.ResetPassRoute + token + utils.TypeRoute + utils.ForgotPassword

		// Determine expired at
		expAt = time.Now().UTC().Add(time.Minute * 5)

		// Import contract send email
		mailContract = mail.New(c.App)
	)

	// Check email and get user data
	userData, err := c.GetUserByEmail(db, ctx, email)
	if err != nil {
		return c.errHandler("model.RequestForgotPassword", err, utils.ErrInvalidEmailPassword)
	}

	// Sending Forgot Password Mail
	err = mailContract.SendMail(mail.UserForgotPassword, mail.MailSubj[mail.UserForgotPassword], email, mail.EmailData{Name: userData.FirstName, Email: email, Link: linkNewPass})
	if err != nil {
		return c.errHandler("model.RequestForgotPassword", err, utils.ErrSendingResetPasswordEmail)
	}

	// Insert verification data into 'verifications' table
	err = c.insertVerificationData(db, ctx, utils.User, utils.ForgotPassword, email, token, false, expAt)
	if err != nil {
		return c.errHandler("model.RequestForgotPassword", err, utils.ErrAddingResetPasswordVerification)
	}

	return nil
}

func (c *Contract) RequestVerifyEmailUser(db *pgxpool.Pool, ctx context.Context, email, types string) error {
	var (
		err error
		// Generate Token 50 digits
		token, _ = utils.Generate(`[a-zA-Z0-9]{50}`)

		// Forgot password route
		link = c.Config.GetString("web_url") + utils.VerifyEmailRoute + token + utils.TypeRoute + types

		// Determine expired at
		expAt = time.Now().UTC().Add(time.Minute * 5)

		// Import contract send email
		mailContract = mail.New(c.App)
	)

	// Check email and get user data
	userData, err := c.GetUserByEmail(db, ctx, email)
	if err != nil {
		return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrInvalidEmailPassword)
	}

	switch types {
	case utils.VerifyRegistration:
		err = mailContract.SendMail(mail.UserVerifyEmail, mail.MailSubj[mail.UserVerifyEmail], email, mail.EmailData{Name: userData.FirstName, Email: email, Link: link})
		if err != nil {
			return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrSendingVerifyEmail)
		}
	case utils.ForgotPassword:
		err = mailContract.SendMail(mail.UserForgotPassword, mail.MailSubj[mail.UserForgotPassword], email, mail.EmailData{Name: userData.FirstName, Email: email, Link: link})
		if err != nil {
			return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrSendingForgotPasswordEmail)
		}
	case utils.UpdateEmail:
		err = mailContract.SendMail(mail.UserUpdateEmail, mail.MailSubj[mail.UserUpdateEmail], email, mail.EmailData{Name: userData.FirstName, Email: email, Link: link})
		if err != nil {
			return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrSendingUpdateEmail)
		}
	default:
		return errors.New(utils.ErrInvalidSendingEmailType)
	}

	err = c.insertVerificationData(db, ctx, utils.User, types, email, token, false, expAt)
	if err != nil {
		return c.errHandler("model.RequestVerifyEmailUser", err, utils.ErrAddingResetPasswordVerification)
	}

	return nil
}

func (c *Contract) ResetPassword(db *pgxpool.Pool, ctx context.Context, userIdentifier, NewPassword, ConfirmPassword string) error {
	var (
		err      error
		dataUser UserEnt
	)
	// Check if new password matches the confirmation
	if NewPassword != ConfirmPassword {
		return errors.New(utils.ErrPasswordMismatch)
	}

	dataUser, err = c.GetUserByUserIdentifier(db, ctx, userIdentifier)
	if err != nil {
		return c.errHandler("model.UpdatePassword", err, utils.ErrFetchingUserPassword)
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(NewPassword), 14)
	if err != nil {
		return c.errHandler("model.UpdatePassword", err, utils.ErrHashingPassword)
	}

	// Update the user's password in the database
	sql := "UPDATE users SET password = $1, updated_date = $3 WHERE id = $2"
	_, err = db.Exec(ctx, sql, string(hashedPassword), dataUser.ID, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.UpdatePassword", err, utils.ErrUpdatingUserPassword)
	}

	return nil
}

func (c *Contract) insertVerificationData(db *pgxpool.Pool, ctx context.Context, actorType, verificationType, email, token string, isUsed bool, expiredDate time.Time) error {
	sql := `INSERT INTO verifications(actor_type, verification_type, email, token, is_used, expired_date, created_date)
        VALUES($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Exec(ctx, sql, actorType, verificationType, email, token, isUsed, expiredDate, time.Now().In(time.UTC))
	if err != nil {
		return err
	}
	return nil
}
