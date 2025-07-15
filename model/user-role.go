package model

type UserRole struct {
	PostgreSQLMetadata
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`

	User User `json:"user,omitempty"`
	Role Role `json:"role,omitempty"`
}

func (UserRole) TableName() string {
	return "user_role"
}
