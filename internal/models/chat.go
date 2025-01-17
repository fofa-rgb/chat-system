package models

import "time"

type Chat struct {
	ID            int64     `db:"id"`
	ApplicationID int64     `db:"application_id"`
	Subject       string    `db:"subject"`
	Number        int64     `db:"number"`
	MessagesCount int64     `db:"messages_count"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
