package database

import (
	"github.com/hafiddna/infrastructure-be-helper/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectToMinio() (client *minio.Client, err error) {
	endpoint := config.Config.App.Minio.Host + ":" + config.Config.App.Minio.Port
	client, err = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(config.Config.App.Minio.AccessKey, config.Config.App.Minio.SecretKey, ""),
		// TODO: Uncomment this line when deploying to production
		//Secure: config.Config.App.Environment == "production",
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
