package model

import (
	"database/sql"
	"time"
)

type Session struct {
	PostgreSQLMetadata
	UserID         sql.NullString `json:"user_id"`
	IpAddress      sql.NullString `json:"ip_address"`
	UserAgent      sql.NullString `json:"user_agent"`
	Payload        string         `json:"payload"`
	LastActivity   time.Time      `json:"last_activity"`
	AppID          string         `json:"app_id"`
	DeviceCategory string         `json:"device_category"`
	DeviceType     string         `json:"device_type"`
	RememberToken  sql.NullString `json:"remember_token"`

	User User `json:"user,omitempty"`
}
