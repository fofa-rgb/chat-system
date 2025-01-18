package models

import "time"

type UserExposedApplication struct {
	Name       string `json:"name" db:"name"`
	Token      string `json:"token" db:"token"`
	ChatsCount int64  `json:"chatsCount" db:"chats_count"`
}

type Application struct {
	Id int64 `db:"id"`
	UserExposedApplication
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
