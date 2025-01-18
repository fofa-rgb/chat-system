package database

import (
	"chat-system/internal/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ChatsDatabaseHandler struct {
	database *sqlx.DB
}

func NewChatsDatabaseHandler() *ChatsDatabaseHandler {
	return &ChatsDatabaseHandler{database: DATABASE}
}

func (r *ChatsDatabaseHandler) InsertChat(appId int64, subject string) (int64, error) {
	var chatNumber int64

	tx := r.database.MustBegin()
	defer tx.Rollback()

	err := tx.Get(&chatNumber, `SELECT COALESCE(MAX(number), 0) FROM Chats WHERE application_id = ? FOR UPDATE`, appId)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch max chat number: %w", err)
	}

	chatNumber++

	_, err = tx.Exec(`INSERT INTO Chats (application_id, subject, number) VALUES (?, ?,?)`, appId, subject, chatNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to insert new chat: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return chatNumber, nil
}
func (r *ChatsDatabaseHandler) GetChatByApplicationIdAndChatNumber(appId int64, chatNumber int64) (models.Chat, error) {
	chat := models.Chat{}
	query := "SELECT * FROM Chats WHERE application_id = ? AND number = ?"
	err := r.database.Get(&chat, query, appId, chatNumber)
	if err != nil {
		return models.Chat{}, fmt.Errorf("failed to get chat: %w", err)
	}
	return chat, nil
}

func (r *ChatsDatabaseHandler) GetAllChatsForAnApp(appId int64) ([]models.Chat, error) {
	allChats := []models.Chat{}
	query := "SELECT * FROM Chats WHERE application_id = ?"
	err := r.database.Select(&allChats, query, appId)
	if err != nil {
		return []models.Chat{}, fmt.Errorf("failed to get chats: %w", err)
	}
	return allChats, nil
}
func (r *ChatsDatabaseHandler) UpdateChatSubject(appId int64, chatNumber int64, newSubject string) (models.Chat, error) {
	chat := models.Chat{}
	query := `
        UPDATE Chats
        SET subject = ?
        WHERE application_id = ? AND number = ?
    `

	// Start a transaction
	tx := r.database.MustBegin()
	defer tx.Rollback()

	_, err := tx.Exec(query, newSubject, appId, chatNumber)
	if err != nil {
		tx.Rollback()
		return models.Chat{}, fmt.Errorf("failed to update chat subject: %w", err)
	}

	fetchQuery := `
        SELECT *
        FROM Chats
        WHERE application_id = ? AND number = ?    `
	err = tx.Get(&chat, fetchQuery, appId, chatNumber)
	if err != nil {
		tx.Rollback()
		return models.Chat{}, fmt.Errorf("failed to fetch updated chat: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return models.Chat{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return chat, nil
}
