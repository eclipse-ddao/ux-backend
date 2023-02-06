package database

import (
	"fmt"
	"log"
	"os"

	"github.com/eclipse-ddao/eclipse-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

func Connect() *gorm.DB {
	PG_HOST := os.Getenv("PG_HOST")
	PG_USERNAME := os.Getenv("PG_USERNAME")
	PG_PASSWORD := os.Getenv("PG_PASSWORD")
	PG_PORT := os.Getenv("PG_PORT")
	DB_NAME := "eclipse"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", PG_HOST, PG_USERNAME, PG_PASSWORD, DB_NAME, PG_PORT)
	dbInstance, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		dbDoesNotExist := strings.Contains(err.Error(), "SQLSTATE 3D000")
		log.Println("DB DOES NOT EXIST?", dbDoesNotExist)
		log.Fatal(err)
	}
	err = dbInstance.AutoMigrate(&models.User{}, &models.Dao{}, &models.File{}, &models.BigFile{}, &models.StorageProvider{}, &models.BigFileProposal{})
	if err != nil {
		log.Fatal(err)
	}
	return dbInstance
}
