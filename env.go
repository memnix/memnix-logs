package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	RabbitMQURL string
	MongoURL    string
)

func LoadVar() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	RabbitMQURL = os.Getenv("RABBITMQ_URL") // Get url from env
	MongoURL = os.Getenv("MONGO_URL")       // Get url from env

}
