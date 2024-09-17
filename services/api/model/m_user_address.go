package model

import (
	"context"
	"database/sql"
	"errors"
	"go-skeleton/lib/utils"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type UserAddressEnt struct {
	ID                int            `db:"id"`
	UserID            int64          `db:"user_id"`
	Title             sql.NullString `db:"title"`
	AddressIdentifier string         `db:"address_identifier"`
	FullAddress       string         `db:"full_address"`
	CreatedDate       time.Time      `db:"created_date"`
	UpdatedDate       sql.NullTime   `db:"updated_date"`
	DeletedDate       sql.NullTime   `db:"deleted_date"`
}

func (c *Contract) GetUserAddressesByUserID(db *pgxpool.Pool, ctx context.Context, userID int64) ([]UserAddressEnt, error) {
	var (
		res []UserAddressEnt

		sql = ` SELECT 
						id, user_id, address_identifier, title, full_address, created_date, updated_date, deleted_date
				FROM user_addresses
				WHERE user_id = $1 AND deleted_date IS NULL
    `
	)
	rows, err := db.Query(ctx, sql, userID)
	if err != nil {
		return res, c.errHandler("model.GetUserAddressesByUserID", err, utils.ErrGettingUserAddresses)
	}
	defer rows.Close()

	for rows.Next() {
		var address UserAddressEnt
		err := rows.Scan(
			&address.ID, &address.UserID, &address.AddressIdentifier, &address.Title, &address.FullAddress,
			&address.CreatedDate, &address.UpdatedDate, &address.DeletedDate,
		)
		if err != nil {
			return res, c.errHandler("model.GetUserAddressesByUserID", err, utils.ErrScanningUserAddresses)
		}
		res = append(res, address)
	}

	if err := rows.Err(); err != nil {
		return res, c.errHandler("model.GetUserAddressesByUserID", err, utils.ErrIteratingUserAddresses)
	}

	return res, nil
}

func (c *Contract) GetAddressByAddressIdentifier(db *pgxpool.Pool, ctx context.Context, addressIdentifier string) (UserAddressEnt, error) {
	var (
		res UserAddressEnt
		sql = `
        SELECT 
				id, user_id, address_identifier, title, full_address, created_date, updated_date, deleted_date
        FROM user_addresses
        WHERE address_identifier = $1 AND deleted_date IS NULL
    `
	)

	err := db.QueryRow(ctx, sql, addressIdentifier).Scan(
		&res.ID, &res.UserID, &res.AddressIdentifier, &res.Title, &res.FullAddress,
		&res.CreatedDate, &res.UpdatedDate, &res.DeletedDate,
	)
	if err != nil {
		return res, c.errHandler("model.GetAddressByAddressIdentifier", err, utils.ErrScanningUserAddresses)
	}

	return res, nil
}

func (c *Contract) InsertUserAddress(db *pgxpool.Pool, ctx context.Context, userID int64, addressIdentifier, title, fullAddress string) error {

	sql := `
        INSERT INTO user_addresses(user_id, address_identifier, title, full_address, created_date)
        VALUES($1, $2, $3, $4, $5)
    `

	_, err := db.Exec(ctx, sql, userID, addressIdentifier, title, fullAddress, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.InsertUserAddress", err, utils.ErrInsertingUserAddress)
	}

	return nil
}

func (c *Contract) UpdateUserAddress(db *pgxpool.Pool, ctx context.Context, userID int64, addressIdentifier, title, fullAddress string) error {
	// Check if the address identifier exists for the given user ID
	var (
		err   error
		count int
	)
	checkSQL := `
	    SELECT COUNT(*)
	    FROM user_addresses
	    WHERE user_id = $1 AND address_identifier = $2
	`
	err = db.QueryRow(ctx, checkSQL, userID, addressIdentifier).Scan(&count)
	if err != nil {
		return c.errHandler("model.UpdateUserAddress", err, utils.ErrCheckingAddressIdentifier)
	}
	if count == 0 {
		return errors.New(utils.EmptyData)
	}

	updateSQL := `
        UPDATE user_addresses
        SET title = $1, full_address = $2,  updated_date = $5
        WHERE user_id = $3 AND address_identifier = $4
    `

	_, err = db.Exec(ctx, updateSQL, title, fullAddress, userID, addressIdentifier, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.UpdateUserAddress", err, utils.ErrUpdatingUserAddress)
	}

	return nil
}

func (c *Contract) DeleteUserAddress(db *pgxpool.Pool, ctx context.Context, userID int64, addressIdentifier string) error {
	// Check if the address identifier exists for the given user ID
	var count int
	checkSQL := `
        SELECT COUNT(*)
        FROM user_addresses
        WHERE user_id = $1 AND address_identifier = $2
    `
	err := db.QueryRow(ctx, checkSQL, userID, addressIdentifier).Scan(&count)
	if err != nil {
		return c.errHandler("model.DeleteUserAddress", err, utils.ErrCheckingAddressIdentifier)
	}
	if count == 0 {
		return errors.New(utils.EmptyData)
	}

	deleteSQL := `
        UPDATE user_addresses SET deleted_date = $3
        WHERE user_id = $1 AND address_identifier = $2
    `

	_, err = db.Exec(ctx, deleteSQL, userID, addressIdentifier, time.Now().UTC())
	if err != nil {
		return c.errHandler("model.DeleteUserAddress", err, utils.ErrDeletingUserAddress)
	}

	return nil
}
