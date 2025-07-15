package model

type UserRole struct {
	PostgreSQLMetadata
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`

	Role       Role       `json:"role,omitempty"`
	Permission Permission `json:"permission,omitempty"`
}

func (UserRole) TableName() string {
	return "user_role"
}
