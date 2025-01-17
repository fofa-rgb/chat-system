package database

import (
	"chat-system/internal/models"
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

func (r *ApplicationsDatabaseHandler) GetApplicationByToken(token string) (models.Application, error) {
	app := models.Application{}
	query := "SELECT * FROM Applications WHERE token = ?"
	err := r.database.Get(&app, query, token)
	if err != nil {
		return models.Application{}, fmt.Errorf("failed to get applications: %w", err)
	}
	return app, nil
}

func (r *ApplicationsDatabaseHandler) GetAllApplications() ([]models.Application, error) {
	allApplications := []models.Application{}
	query := "SELECT * FROM Applications"
	err := r.database.Select(&allApplications, query)
	if err != nil {
		return []models.Application{}, fmt.Errorf("failed to get applications: %w", err)
	}
	return allApplications, nil
}

func (r *ApplicationsDatabaseHandler) UpdateApplicationName(token string, name string) error {
	query := "UPDATE Applications SET name = ? WHERE token = ?"
	_, err := r.database.Exec(query, name, token)
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	return nil
}
