package database

import (
	"fmt"
	"os"

	"github.com/AlexWilliam12/silent-signal/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ExecMigration() {
	db := OpenConn()
	if err := db.AutoMigrate(&models.User{}, &models.PrivateMessage{}, &models.Group{}, &models.GroupMessage{}); err != nil {
		panic(fmt.Errorf("failed to execute migration: %v", err))
	}
}

func OpenConn() *gorm.DB {
	db, err := gorm.Open(getConn(), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func getConn() gorm.Dialector {
	return postgres.Open(fmt.Sprintf(`
	host=%s
	user=%s
	password=%s
	dbname=%s
	port=%s
	sslmode=disable
	TimeZone=America/Sao_Paulo`,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT")))
}
