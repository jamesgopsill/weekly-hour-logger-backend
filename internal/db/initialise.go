package db

import (
	"jamesgopsill/resource-logger-backend/internal/config"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var Connection *gorm.DB

func Initialise() {
	var err error
	log.Info().Msg("Connecting to database")
	Connection, err = gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Info().Msg("Connected to database")
	err = Connection.AutoMigrate(&User{}, &Group{}, &Resource{})
	if err != nil {
		panic(err)
	}
}
