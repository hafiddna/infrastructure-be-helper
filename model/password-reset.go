package model

import (
	"database/sql"
	"time"
)

type PasswordReset struct {
	PostgreSQLMetadata
	UserID    string       `json:"user_id"`
	Type      string       `json:"type"`
	Token     string       `json:"token"`
	Attempts  int          `json:"attempts"`
	ExpiresAt time.Time    `json:"expires_at"`
	UsedAt    sql.NullTime `json:"used_at"`

	User User `json:"user,omitempty"`
}
