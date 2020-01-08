package database

import (
	"botBoilerplate/modules"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var db *gorm.DB

func Init(file string) error {
	log.Printf("Connecting to database...")

	tdb, err := gorm.Open("sqlite3", file)
	if err != nil {
		return err
	}

	db = tdb

	log.Print("Creating tables...")
	models := modules.GetModels()
	for _, model := range models {
		log.Printf("Creating %s.", model.Name)

		if db.HasTable(model.Model) {
			db.AutoMigrate(model.Model)
		} else {
			db.CreateTable(model.Model)
		}

		if db.Error != nil {
			return db.Error
		}
	}

	return err
}

func GetDB() *gorm.DB {
	return db
}

func Close() error {
	return db.Close()
}
