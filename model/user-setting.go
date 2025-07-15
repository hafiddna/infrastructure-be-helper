package model

type UserSetting struct {
	MongoDBMetadata
	UserID string `json:"user_id" bson:"user_id"`
}
