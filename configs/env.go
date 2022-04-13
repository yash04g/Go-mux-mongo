package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load() // Loading the .env file
	if err != nil {
		log.Fatal("Error in loading .env file")
	}

	return os.Getenv("MONGOURI") // Returning the mongo string
}
