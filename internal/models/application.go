package models

import "time"

type UserExposedApplication struct {
	Name       string `json:"name" db:"name"`
	Token      string `json:"token" db:"token"`
	ChatsCount int    `json:"chatsCount" db:"chats_count"`
}

type Application struct {
	ID int64 `db:"id"`
	UserExposedApplication
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
