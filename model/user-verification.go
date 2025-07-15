package model

import (
	"database/sql"
	"time"
)

type UserVerification struct {
	PostgreSQLMetadata
	UserID     string       `json:"user_id"`
	Type       string       `json:"type"`
	Token      string       `json:"token"`
	Attempts   int          `json:"attempts"`
	ExpiresAt  time.Time    `json:"expires_at"`
	VerifiedAt sql.NullTime `json:"verified_at"`

	User User `json:"user,omitempty"`
}
