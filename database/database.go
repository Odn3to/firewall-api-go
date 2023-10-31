package database

import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"

	"firewall-api-go/models"
)
var (
	DB  *gorm.DB
)

func ConectaNoBD() {
	db, err := gorm.Open(sqlite.Open("config_firewall.db"), &gorm.Config{})
  	if err != nil {
    	panic("failed to connect database")
  	}

	DB = db

	db.AutoMigrate(&models.Regra{})
}

func GetDataBase() *gorm.DB {
	return DB
}