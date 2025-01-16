package database

import "github.com/jmoiron/sqlx"

type messagesDatabaseHandler struct {
	database *sqlx.DB
}

func (r *messagesDatabaseHandler) InsertMessage(chatNumber int64, messageNumber int, chatId int64, body string) error {
	query := `INSERT INTO Messages (chat_id, number, body)
		      VALUES (?, ?, ?)`
	_, err := r.database.Exec(query, chatId, messageNumber, body)
	if err != nil {
		return err
	}
	return nil
}
func (r *messagesDatabaseHandler) UpdateMessage(chatId int64, messageNumber int, newBody string) error {
	query := `UPDATE Messages
			  SET body = ?,
			  WHERE chat_id = ? AND number = ?`
	_, err := r.database.Exec(query, newBody, chatId, messageNumber)
	if err != nil {
		return err
	}
	return nil
}
