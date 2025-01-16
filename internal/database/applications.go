package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ApplicationsDatabaseHandler struct {
	database *sqlx.DB
}

func NewApplicationsDatabaseHandler() *ApplicationsDatabaseHandler {
	return &ApplicationsDatabaseHandler{database: DATABASE}
}

func (r *ApplicationsDatabaseHandler) InsertApplication(name string, token string) error {
	query := `
        INSERT INTO Applications (name, token)
        VALUES (?, ?)
    `

	_, err := r.database.Exec(query, name, token)
	if err != nil {
		return fmt.Errorf("failed to insert application: %w", err)
	}

	return nil
}

func (r *ApplicationsDatabaseHandler) UpdateApplication(id int64, name string) error {
	query := "UPDATE Applications SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2"
	_, err := r.database.Exec(query, name, id)
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	return nil
}
