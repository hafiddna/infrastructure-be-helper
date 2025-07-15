package model

type UserProfile struct {
	MongoDBMetadata
	UserID   string `json:"user_id" bson:"user_id"`
	FullName string `json:"full_name" bson:"full_name,omitempty"`
	NickName string `json:"nick_name" bson:"nick_name,omitempty"`
}
