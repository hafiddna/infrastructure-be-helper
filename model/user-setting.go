package model

import "go.mongodb.org/mongo-driver/v2/bson"

type UserSetting struct {
	ID        bson.ObjectID          `json:"id" bson:"_id,omitempty"`
	UserID    string                 `json:"user_id" bson:"user_id"`
	CreatedAt bson.DateTime          `json:"created_at" bson:"created_at"`
	UpdatedAt bson.DateTime          `json:"updated_at" bson:"updated_at"`
	DeletedAt *bson.DateTime         `json:"deleted_at" bson:"deleted_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata" bson:"metadata,omitempty"`
}
