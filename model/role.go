package model

import "database/sql"

type Role struct {
	PostgreSQLMetadata
	Name        string         `json:"name"`
	DisplayName string         `json:"display_name"`
	Description sql.NullString `json:"description"`
	Category    string         `json:"category"`

	Users           []User           `json:"users,omitempty" gorm:"many2many:user_role"`
	UserRoles       []UserRole       `json:"user_roles,omitempty"`
	Permissions     []Permission     `json:"permissions,omitempty" gorm:"many2many:role_permission"`
	RolePermissions []RolePermission `json:"role_permissions,omitempty"`
}
