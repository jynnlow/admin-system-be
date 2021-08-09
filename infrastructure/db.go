package infrastructure

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"admin-system-be/models"
)

// ConnectDB ...
func ConnectDB() (*gorm.DB, error) {
	//to connect mysql db - username:password@protocol(address)/dbname?param=value
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = dbConn.AutoMigrate(&models.User{}, &models.Secret{})
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
