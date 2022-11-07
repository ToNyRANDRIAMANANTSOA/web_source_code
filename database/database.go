package database

import (
	"back/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	DB *gorm.DB
}

var Database DB

func Connect() {
	db, err := gorm.Open(sqlite.Open("student.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	Database = DB{
		DB: db,
	}
	if err != nil {
		log.Fatal(err)
		return
	}
	Database.DB.AutoMigrate(&models.Student{}, &models.EC{}, &models.Note{})
	// studentTest := models.Student{
	// 	Nom:     "Nomentsoa Rakotonirina",
	// 	Adresse: "Beravina",
	// 	Sexe:    "M",
	// 	Niveau:  "L2",
	// 	Annee:   2021,
	// }
	// Database.DB.Create(&studentTest)
}
