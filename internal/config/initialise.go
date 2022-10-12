package config

import (
	"os"

	"github.com/rs/zerolog/log"
)

var DBPath string
var Secret string
var MySigningKey []byte
var Issuer string

func Initalise(dev *bool) {

	if *dev {
		log.Info().Msg("Starting in development mode.")
		Secret = "shhh"
		Issuer = "www.test.com"
		DBPath = "data/test.db"
		return
	}

	Secret := os.Getenv("GO_REST_JWT_SECRET")
	if Secret == "" {
		panic("No secret.")
	}
	MySigningKey = []byte(Secret)

	Issuer = os.Getenv("GO_REST_JWT_ISSUER")
	if Issuer == "" {
		panic("No issuer.")
	}

	DBPath = os.Getenv("GO_REST_DB_PATH")
	if DBPath == "" {
		panic("No DB Path.")
	}

}
