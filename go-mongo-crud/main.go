package main

import (
	"fmt"
	"log"

	"github.com/go-snippets/go-mongo-crud/app"
	"github.com/go-snippets/go-mongo-crud/router"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	app.SetupConn()
	router.InitServer()
}
