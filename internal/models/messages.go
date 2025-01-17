package models

import "time"

type Message struct {
	ID        int64     `db:"id"`
	ChatID    int64     `db:"chat_id"`
	Number    int64     `db:"number"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
