package model

import "database/sql"

type User struct {
	PostgreSQLMetadata
	Email                  sql.NullString `json:"email"`
	EmailVerifiedAt        sql.NullTime   `json:"email_verified_at"`
	Phone                  sql.NullString `json:"phone"`
	PhoneVerifiedAt        sql.NullTime   `json:"phone_verified_at"`
	Username               sql.NullString `json:"username"`
	Password               string         `json:"password"`
	Pin                    sql.NullString `json:"pin"`
	TwoFactorSecret        sql.NullString `json:"two_factor_secret"`
	TwoFactorRecoveryCodes sql.NullString `json:"two_factor_recovery_codes"`
	TwoFactorConfirmedAt   sql.NullTime   `json:"two_factor_confirmed_at"`
	Status                 string         `json:"status"`

	PasswordResets []PasswordReset    `json:"password_resets,omitempty"`
	Verifications  []UserVerification `json:"verifications,omitempty"`
	Sessions       []Session          `json:"sessions,omitempty"`
	UserRoles      []UserRole         `json:"user_roles,omitempty"`
	Roles          []Role             `json:"roles,omitempty" gorm:"many2many:user_role"`
}
