package database

import (
	"fmt"
	"github.com/hafiddna/infrastructure-be-helper/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectToMongoDB() (db *mongo.Database, err error) {
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		config.Config.App.MongoDB.Username,
		config.Config.App.MongoDB.Password,
		config.Config.App.MongoDB.Host,
		config.Config.App.MongoDB.Port,
		config.Config.App.MongoDB.Database,
	)

	client, err := mongo.Connect(options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	err = client.Ping(nil, nil)
	if err != nil {
		return nil, err
	}

	db = client.Database(config.Config.App.MongoDB.Database)
	return db, nil
}
