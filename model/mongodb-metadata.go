package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoDBMetadata struct {
	ID        bson.ObjectID          `json:"id" bson:"_id,omitempty"`
	CreatedAt bson.Timestamp         `json:"created_at" bson:"created_at"`
	UpdatedAt bson.Timestamp         `json:"updated_at" bson:"updated_at"`
	DeletedAt bson.Timestamp         `json:"deleted_at" bson:"deleted_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata" bson:"metadata"`
}

func (m *MongoDBMetadata) setMetadata(metadata map[string]interface{}) {
	if m.Metadata == nil {
		m.Metadata = make(map[string]interface{})
	}

	for key, value := range metadata {
		m.Metadata[key] = value
	}
}

func (m *MongoDBMetadata) Created(userID string) {
	initialMetadata := map[string]interface{}{
		"created_by": userID,
		"updated_by": userID,
	}
	m.setMetadata(initialMetadata)
}

func (m *MongoDBMetadata) Updated(userID string) {
	updatedMetadata := map[string]interface{}{
		"updated_by": userID,
	}
	m.setMetadata(updatedMetadata)
}

func (m *MongoDBMetadata) Deleted(userID string) {
	deletedMetadata := map[string]interface{}{
		"deleted_by": userID,
	}
	m.setMetadata(deletedMetadata)
}
