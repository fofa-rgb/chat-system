package database

import (
	"chat-system/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MessagesDatabaseHandler struct {
	database *sqlx.DB
}

func NewMessagesDatabaseHandler() *MessagesDatabaseHandler {
	return &MessagesDatabaseHandler{database: DATABASE}
}

func (r *MessagesDatabaseHandler) InsertMessage(chatId int64, body string) (int64, error) {
	var messageNumber int64

	tx := r.database.MustBegin()
	defer tx.Rollback()

	err := tx.Get(&messageNumber, `SELECT COALESCE(MAX(number), 0) FROM Messages WHERE chat_id = ? FOR UPDATE`, chatId)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch max message number: %w", err)
	}

	messageNumber++

	_, err = tx.Exec(`INSERT INTO Messages (chat_id, body, number) VALUES (?, ?,?)`, chatId, body, messageNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new message: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return messageNumber, nil
}
func (r *MessagesDatabaseHandler) GetMessageByChatIdAndMessageNumber(chatId int64, messageNumber int64) (models.Message, error) {
	message := models.Message{}
	query := "SELECT * FROM Messages WHERE chat_id = ? AND number = ?"
	err := r.database.Get(&message, query, chatId, messageNumber)
	if err != nil {
		return models.Message{}, fmt.Errorf("failed to get message: %w", err)
	}
	return message, nil
}

func (r *MessagesDatabaseHandler) GetAllMessagesForAChat(chatId int64) ([]models.Message, error) {
	allMessages := []models.Message{}
	query := "SELECT * FROM Messages WHERE chat_id = ?"
	err := r.database.Select(&allMessages, query, chatId)
	if err != nil {
		return []models.Message{}, fmt.Errorf("failed to get messages: %w", err)
	}
	return allMessages, nil
}

func (r *MessagesDatabaseHandler) UpdateMessageBody(chatId int64, messageNumber int64, newBody string) (models.Message, error) {
	updatedMessage := models.Message{}
	query := `
        UPDATE Messages
        SET body = ?
        WHERE chat_id = ? AND number = ?
    `

	tx := r.database.MustBegin()
	defer tx.Rollback()

	_, err := tx.Exec(query, newBody, chatId, messageNumber)
	if err != nil {
		tx.Rollback()
		return models.Message{}, fmt.Errorf("failed to update message subject: %w", err)
	}

	fetchQuery := `
        SELECT *
        FROM Messages
        WHERE chat_id = ? AND number = ?    
		`
	err = tx.Get(&updatedMessage, fetchQuery, chatId, messageNumber)
	if err != nil {
		tx.Rollback()
		return models.Message{}, fmt.Errorf("failed to fetch updated chat: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return models.Message{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedMessage, nil
}
