package repository

import (
	"fmt"
	"log"
	"news_service/internal/logger"
	"news_service/internal/model"
	_ "news_service/internal/model"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		return nil, fmt.Errorf("there isn't dsn")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db
	err = migrateModels(db)
	if err != nil {
		return nil, err
	}
	logger.ZapLogger.Info("database is ready")
	return db, nil
}

func migrateModels(db *gorm.DB) error {
	err := db.AutoMigrate(&model.UserModel{}, &model.NewsModel{}, &model.BlockModel{}, &model.ImageModel{})
	if err != nil {
		return err
	}
	log.Println("Tables are ok")
	return nil
}