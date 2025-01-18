package models

import "time"

type UserExposedChat struct {
	Subject       string `db:"subject"`
	Number        int64  `db:"number"`
	MessagesCount int64  `db:"messages_count"`
}
type Chat struct {
	Id            int64 `db:"id"`
	ApplicationId int64 `db:"application_id"`
	UserExposedChat
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
