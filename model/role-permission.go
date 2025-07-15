package model

type RolePermission struct {
	PostgreSQLMetadata
	PermissionID string `json:"permission_id"`
	RoleID       string `json:"role_id"`

	Permission Permission `json:"permission,omitempty"`
	Role       Role       `json:"role,omitempty"`
}

func (RolePermission) TableName() string {
	return "role_permission"
}
