package models

import "time"

type Application struct {
	ID         int64     `db:"id"`
	Name       string    `db:"name"`
	Token      string    `db:"token"`
	ChatsCount int       `db:"chats_count"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
