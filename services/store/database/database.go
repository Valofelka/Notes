package database

import (
	"notes_project/services/store/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(conf *config.Config) (*gorm.DB, error) {
	dsn := "user=" + conf.DB.User + " password=" + conf.DB.Password + " dbname=" + conf.DB.DBName + " host=" + conf.DB.Host + " sslmode=" + conf.DB.SSLMode

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	return db, nil

}
