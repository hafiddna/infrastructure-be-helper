package model

import "database/sql"

type Permission struct {
	PostgreSQLMetadata
	Name        string         `json:"name"`
	DisplayName string         `json:"display_name"`
	Description sql.NullString `json:"description"`
	Module      sql.NullString `json:"module"`
	Category    string         `json:"category"`

	Roles           []Role           `json:"roles,omitempty" gorm:"many2many:role_permission"`
	RolePermissions []RolePermission `json:"role_permissions,omitempty"`
}
