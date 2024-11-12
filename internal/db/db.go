package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"viabl.ventures/gossr/internal/config"
	"viabl.ventures/gossr/internal/db/models"
)

func InitDB(conf *config.EnvVars) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		conf.DbHost,
		conf.DbUser,
		conf.DbPassword,
		conf.DbName,
		conf.DbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("Connected Successfully to Database")
	var DB *gorm.DB = db

	// Auto Migrate the schemas
	DB.AutoMigrate(&models.AdminUser{}, &models.LoginCode{})

	return DB
}
