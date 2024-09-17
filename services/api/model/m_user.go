package model

import (
	"context"
	"database/sql"
	"errors"
	"go-skeleton/lib/utils"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserEnt struct {
	ID             int            `db:"id"`
	UserIdentifier string         `db:"user_identifier"`
	FirstName      string         `db:"first_name"`
	LastName       sql.NullString `db:"last_name"`
	Email          string         `db:"email"`
	AvatarURL      sql.NullString `db:"avatar_url"`
	Description    sql.NullString `db:"description"`
	Password       string         `db:"password"`
	IsVerified     bool           `db:"is_verify"`
	CreatedDate    time.Time      `db:"created_date"`
	UpdatedDate    sql.NullTime   `db:"updated_date"`
	DeletedDate    sql.NullTime   `db:"deleted_date"`
}

func (c *Contract) GetUserByEmail(db *pgxpool.Pool, ctx context.Context, email string) (UserEnt, error) {
	var res UserEnt

	sql := `SELECT id, user_identifier, first_name, last_name, email, avatar_url, description, password, is_verify, created_date, updated_date, deleted_date
            FROM users
            WHERE email = $1 AND deleted_date IS NULL`

	err := db.QueryRow(ctx, sql, email).Scan(
		&res.ID,
		&res.UserIdentifier,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.AvatarURL,
		&res.Description,
		&res.Password,
		&res.IsVerified,
		&res.CreatedDate,
		&res.UpdatedDate,
		&res.DeletedDate,
	)

	if err != nil {
		return res, c.errHandler("model.GetUserByEmail", err, utils.ErrGettingUserByEmail)
	}

	return res, nil
}

func (c *Contract) GetUserByUserIdentifier(db *pgxpool.Pool, ctx context.Context, userIdentifier string) (UserEnt, error) {
	var res UserEnt

	sql := `SELECT id, user_identifier, first_name, last_name, email, avatar_url, description, password, is_verify, created_date, updated_date, deleted_date
            FROM users
            WHERE user_identifier = $1 AND deleted_date IS NULL`

	err := db.QueryRow(ctx, sql, userIdentifier).Scan(
		&res.ID,
		&res.UserIdentifier,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.AvatarURL,
		&res.Description,
		&res.Password,
		&res.IsVerified,
		&res.CreatedDate,
		&res.UpdatedDate,
		&res.DeletedDate,
	)

	if err != nil {
		return res, c.errHandler("model.GetUserByUserIdentifier", err, utils.ErrRetrievingUserByUserIdentifier)
	}

	return res, nil
}

func (c *Contract) UpdateUserProfile(db *pgxpool.Pool, ctx context.Context, userIdentifier, firstName, lastName, description, avatarURL string) error {
	sql := `
		UPDATE users
		SET avatar_url = $1, first_name = $2, last_name = $3, description = $4 
		WHERE user_identifier = $5
	`

	_, err := db.Exec(ctx, sql, avatarURL, firstName, lastName, description, userIdentifier)
	if err != nil {
		return c.errHandler("model.UpdateUserProfile", err, utils.ErrUpdatingUserProfile)
	}

	return nil
}

func (c *Contract) UpdateUserEmail(db *pgxpool.Pool, ctx context.Context, userIdentifier, email string) error {
	sql := `
		UPDATE users
		SET email = $1
		WHERE user_identifier = $2
	`

	_, err := db.Exec(ctx, sql, email, userIdentifier)
	if err != nil {
		return c.errHandler("model.UpdateUserEmail", err, utils.ErrUpdatingUserEmail)
	}

	return nil
}

func (c *Contract) UpdatePasswordUser(db *pgxpool.Pool, ctx context.Context, userIdentifier, OldPassword, NewPassword, ConfirmPassword string) error {
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
		return c.errHandler("model.UpdatePasswordUser", err, utils.ErrFetchingUserPassword)
	}

	// Validate old password
	err = bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(OldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.errHandler("model.UpdatePasswordUser", err, utils.ErrHashingPassword)
	}

	// Update the user's password in the database
	sql := "UPDATE users SET password = $1, updated_date = $3 WHERE id = $2"
	_, err = db.Exec(ctx, sql, string(hashedPassword), dataUser.ID, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.UpdatePasswordUser", err, utils.ErrUpdatingUserPassword)
	}

	return nil
}
