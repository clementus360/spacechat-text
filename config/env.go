package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err!=nil {
		fmt.Println("Failed to load env")
		log.Fatal(err)
	}
}
