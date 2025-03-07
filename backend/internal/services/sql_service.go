package services

import (
	"fmt"
	"log"
	"os"

	"github.com/mateuse/yahoo-fantasy-analyzer/internal/repositories"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToMySQL() error {
	sqlUser := os.Getenv("SQL_USER")
	sqlPassword := os.Getenv("SQL_PASSWORD")
	sqlHost := os.Getenv("SQL_HOST")
	sqlPort := os.Getenv("SQL_PORT")
	sqlDbName := os.Getenv("SQL_DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		sqlUser, sqlPassword, sqlHost, sqlPort, sqlDbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to mysql: %w", err)
	}

	repositories.Init(DB)

	log.Println("Connected to MySQL successfully")
	return nil
}
