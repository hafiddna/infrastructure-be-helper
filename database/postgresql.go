package database

import (
	"fmt"
	"github.com/hafiddna/infrastructure-be-helper/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func ConnectToPostgreSQL() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.Config.App.PostgreSQL.Host,
		config.Config.App.PostgreSQL.Username,
		config.Config.App.PostgreSQL.Password,
		config.Config.App.PostgreSQL.Database,
		config.Config.App.PostgreSQL.Port,
	)

	var gormLogger logger.Interface
	if config.Config.App.Environment == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		//DryRun: true,
		//SkipHooks: true, // should be using &gorm.Session{SkipHooks: true}
		//QueryFields: true,
		//SkipDefaultTransaction: true,
		//PrepareStmt:            true,
		//DisableAutomaticPing:   true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)

	return db, nil
}
