package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Entries EntryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Entries: EntryModel{DB: db},
	}
}
