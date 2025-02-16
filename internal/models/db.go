package models

import "github.com/jackc/pgx/v5"

type EmptyRow struct{}

func (r EmptyRow) Scan(...interface{}) error {
	return pgx.ErrNoRows
}
