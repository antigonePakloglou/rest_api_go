// initialisation à la bdd ici
package models

import (
	"rest/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	//ouverture connection bdd
	database, err := gorm.Open(postgres.Open(env.DSN), &gorm.Config{})

	if err != nil {
		panic("Echec de la connection à la bdd !")
	}
	//migartion  model
	err = database.AutoMigrate(&Users{}, &Articles{}, &Comments{})

	if err != nil {
		return
	}

	DB = database
}
