package models

import "time"

type Chat struct {
	ID            int64     `db:"id"`
	ApplicationID int64     `db:"application_id"`
	Number        int       `db:"number"`
	MessagesCount int       `db:"messages_count"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
