package utils

import (
	"os"

	"github.com/joho/godotenv"
)


func IsSecure()bool{
	if err := godotenv.Load();err != nil {
		panic("Error loading .env file")
	}
	if os.Getenv("ENVIRONMENT") == "DEVELOPMENT" {
		return false
	} else {
		return true
	}
}





