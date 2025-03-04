package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Overload(); err != nil {
		log.Println("Warning: .env file is missing. Using environment variables")
	}
	Connect()
}

func Initialize() error {
	if err := Connect(); err != nil {
		return fmt.Errorf("%s", err)
	}

	if err := ConnectSession(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
