package model

import (
	"github.com/hafiddna/infrastructure-be-helper/helper"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type PostgreSQLMetadata struct {
	ID        string         `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Metadata  datatypes.JSON `json:"metadata"`
}

func (m *PostgreSQLMetadata) setMetadata(metadata map[string]interface{}) {
	var oldMetadata map[string]interface{}
	_ = helper.JSONUnmarshal(m.Metadata, &oldMetadata)

	for key, value := range metadata {
		oldMetadata[key] = value
	}

	stringMetadata := helper.JSONMarshal(oldMetadata)
	m.Metadata = datatypes.JSON(stringMetadata)
}

func (m *PostgreSQLMetadata) Created(userID *string) {
	initialMetadata := map[string]interface{}{
		"created_by": userID,
		"updated_by": userID,
	}
	m.setMetadata(initialMetadata)
}

func (m *PostgreSQLMetadata) Updated(userID *string) {
	updatedMetadata := map[string]interface{}{
		"updated_by": userID,
	}
	m.setMetadata(updatedMetadata)
}

func (m *PostgreSQLMetadata) Deleted(userID *string) {
	deletedMetadata := map[string]interface{}{
		"deleted_by": userID,
	}
	m.setMetadata(deletedMetadata)
}
