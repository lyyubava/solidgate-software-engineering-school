package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Database(connString string) {
	database, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		panic("Fail to connect to db")
	}
	DB = database
}
