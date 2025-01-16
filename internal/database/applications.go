package database

import "github.com/jmoiron/sqlx"

type applicationsDatabaseHandler struct {
	database *sqlx.DB
}

func (r *applicationsDatabaseHandler) InsertApplication() error {

	return nil
}
func (r *applicationsDatabaseHandler) UpdateApplication() error {

	return nil
}
