package data

import (
	"database/sql"
	"time"
)

type Entry struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Rating      int8      `json:"rating"`
	Testimonial string    `json:"testimonial"`
	LastName    string    `json:"lastname"`
	FirstName   string    `json:"firstname"`
	Occupation  string    `json:"occupation"`
	Github      string    `json:"github"`
}

type EntryModel struct {
	DB *sql.DB
}

func (e EntryModel) Insert(entry *Entry) error {
	query := `
	INSERT INTO entries (rating, testimonial, last_name, first_name, occupation, github)
	VALUES ($1, $2, $3, $4, $5, $6)`

	args := []any{entry.Rating, entry.Testimonial, entry.LastName, entry.FirstName, entry.Occupation, entry.Github}

	_, err := e.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (e EntryModel) GetAll() ([]*Entry, error) {
	query := `
	SELECT id, created_at, rating, testimonial, last_name, first_name, occupation, github FROM entries`

	rows, err := e.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	entries := []*Entry{}

	for rows.Next() {
		entry := &Entry{}
		err = rows.Scan(&entry.ID, &entry.CreatedAt, &entry.Rating, &entry.Testimonial, &entry.LastName, &entry.FirstName, &entry.Occupation, &entry.Github)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
