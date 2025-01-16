package database

import "github.com/jmoiron/sqlx"

type chatsDatabaseHandler struct {
	database *sqlx.DB
}

func (r *chatsDatabaseHandler) GetChatIdFromAppTokenAndChatNum(appToken string, chatNumber int64) (int64, error) {
	var chatID int64
	query := `
		SELECT c.id
		FROM Chats c
		INNER JOIN Applications a ON c.application_id = a.id
		WHERE a.token = ? AND c.number = ?
	`
	err := r.database.Get(&chatID, query, appToken, chatNumber)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}
