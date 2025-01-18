package models

import "time"

type UserExposedMessage struct {
	Number int64  `db:"number"`
	Body   string `db:"body"`
}

type Message struct {
	Id     int64 `db:"id"`
	ChatId int64 `db:"chat_id"`
	UserExposedMessage
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
