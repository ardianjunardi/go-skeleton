package model

import (
	"context"
	"database/sql"
	"fmt"
	"go-skeleton/lib/utils"
	"go-skeleton/services/api/request"
	"math"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SettingEnt struct {
	Id           int64        `db:"id"`
	SettingCode  string       `db:"setting_code"`
	SetGroup     string       `db:"set_group"`
	SetKey       string       `db:"set_key"`
	SetLabel     string       `db:"set_label"`
	SetOrder     int          `db:"set_order"`
	ContentType  string       `db:"content_type"`
	ContentValue string       `db:"content_value"`
	IsActive     bool         `db:"is_active"`
	CreatedDate  time.Time    `db:"created_date"`
	UpdatedDate  sql.NullTime `db:"updated_date"`
}

var setType = []string{"json_arr", "json_obj", "bool", "string"}

func (c *Contract) GetSetting(db *pgxpool.Pool, ctx context.Context, param request.SettingParam) ([]SettingEnt, error) {
	var (
		err        error
		list       []SettingEnt
		where      []string
		paramQuery []interface{}
		totalData  int

		query = `SELECT 
		id, setting_code, set_group, set_key, set_label, set_order, content_type, content_value, is_active 
		FROM settings`
	)

	// Populate Search
	if len(param.Keyword) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, "%"+param.Keyword+"%")
		orWhere = append(orWhere, fmt.Sprintf("set_label iLIKE $%d", len(paramQuery)))
		orWhere = append(orWhere, fmt.Sprintf("content_value iLIKE $%d", len(paramQuery)))
		where = append(where, "("+strings.Join(orWhere, " OR ")+")")
	}
	if len(param.IsActive) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.IsActive)
		orWhere = append(orWhere, fmt.Sprintf("is_active = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}
	if len(param.SetGroup) > 0 {
		var orWhere []string
		paramQuery = append(paramQuery, param.SetGroup)
		orWhere = append(orWhere, fmt.Sprintf("set_group = $%d", len(paramQuery)))
		where = append(where, strings.Join(orWhere, " AND "))
	}

	// Append All Where Conditions
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	{
		newQcount := `SELECT COUNT(*) FROM ( ` + query + ` ) AS data`
		err := db.QueryRow(ctx, newQcount, paramQuery...).Scan(&totalData)
		if err != nil {
			return list, c.errHandler("model.GetSetting", err, utils.ErrCountingListSetting)
		}
		param.Count = totalData
	}

	// Select Max Page
	if param.Count > param.Limit && param.Page > int(param.Count/param.Limit) {
		param.Page = int(math.Ceil(float64(param.Count) / float64(param.Limit)))
	}

	// Limit and Offset
	param.Offset = (param.Page - 1) * param.Limit
	query += " ORDER BY " + param.Order + " " + param.Sort + " "

	paramQuery = append(paramQuery, param.Offset)
	query += fmt.Sprintf("offset $%d ", len(paramQuery))

	paramQuery = append(paramQuery, param.Limit)
	query += fmt.Sprintf("limit $%d ", len(paramQuery))

	rows, err := db.Query(ctx, query, paramQuery...)
	if err != nil {
		return list, c.errHandler("model.GetSetting", err, utils.ErrGettingListSetting)
	}

	defer rows.Close()
	for rows.Next() {
		var data SettingEnt
		err = rows.Scan(&data.Id, &data.SettingCode, &data.SetGroup, &data.SetKey, &data.SetLabel, &data.SetOrder, &data.ContentType, &data.ContentValue, &data.IsActive)
		if err != nil {
			return list, c.errHandler("model.GetSetting", err, utils.ErrScanningListSetting)
		}
		list = append(list, data)
	}
	return list, nil
}

func (c *Contract) GetSettingByCode(db *pgxpool.Pool, ctx context.Context, code string) (SettingEnt, error) {
	var (
		err  error
		data SettingEnt
		sql  = `SELECT id, setting_code, set_group, set_key, set_label, set_order, content_type, content_value, is_active 
		FROM  settings 
		WHERE setting_code = $1`
	)
	err = db.QueryRow(ctx, sql, code).Scan(&data.Id, &data.SettingCode, &data.SetGroup, &data.SetKey, &data.SetLabel, &data.SetOrder, &data.ContentType, &data.ContentValue, &data.IsActive)
	if err != nil {
		return data, c.errHandler("model.GetSettingByCode", err, utils.ErrGettingSettingByCode)
	}

	return data, nil
}

func (c *Contract) GetSettingValueByKey(db *pgxpool.Pool, ctx context.Context, key string) (string, error) {
	var (
		err error
		res string
		sql = `SELECT content_value 
		FROM  settings 
		WHERE set_key = $1 AND is_active = true`
	)
	err = db.QueryRow(ctx, sql, key).Scan(&res)
	if err != nil {
		return res, c.errHandler("model.GetSettingValueByKey", err, utils.ErrGettingSettingByKey)
	}

	return res, nil
}

func (c *Contract) AddSetting(db *pgxpool.Pool, ctx context.Context, code, group, label string, order int, contentType, content string, isActive bool) error {
	if !utils.Contains(setType, contentType) {
		return fmt.Errorf("%s", "wrong content type value for settings(json_arr|json_obj|bool|string")
	}

	sql := `INSERT INTO settings(setting_code, set_group, set_key, set_label, set_order, content_type, content_value, is_active, created_date)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	_, err := db.Exec(ctx, sql, code, group, content, label, order, contentType, content, isActive, time.Now().In(time.UTC))
	if err != nil {
		return c.errHandler("model.AddSetting", err, utils.ErrAddingSetting)
	}

	return nil
}
func (c *Contract) UpdateSetting(db *pgxpool.Pool, ctx context.Context, code, label string, order int, content string, isActive bool) error {
	var (
		err error
		sql = `
		UPDATE settings 
		SET set_key=$1,set_label=$2,set_order=$3,content_value=$4,is_active=$5,updated_date=$6
		WHERE setting_code=$7`
	)
	_, err = db.Exec(ctx, sql, content, label, order, content, isActive, time.Now().In(time.UTC), code)
	if err != nil {
		return c.errHandler("model.UpdateSetting", err, utils.ErrUpdatingSetting)
	}

	return nil
}
