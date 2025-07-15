package model

import "go.mongodb.org/mongo-driver/v2/bson"

type UserProfile struct {
	ID        bson.ObjectID          `json:"id" bson:"_id,omitempty"`
	UserID    string                 `json:"user_id" bson:"user_id"`
	FullName  string                 `json:"full_name" bson:"full_name,omitempty"`
	NickName  string                 `json:"nick_name" bson:"nick_name,omitempty"`
	CreatedAt bson.DateTime          `json:"created_at" bson:"created_at"`
	UpdatedAt bson.DateTime          `json:"updated_at" bson:"updated_at"`
	DeletedAt *bson.DateTime         `json:"deleted_at" bson:"deleted_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata" bson:"metadata,omitempty"`
}
