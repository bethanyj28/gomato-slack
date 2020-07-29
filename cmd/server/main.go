package main

import (
	"log"

	"github.com/bethanyj28/gomato-slack/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("environment.env"); err != nil {
		log.Fatal("failed to load env")
	}

	server := routes.NewServer()

	server.BuildRoutes()

	if err := server.Router.Run(":8080"); err != nil {
		log.Fatal("server quit unexpectedly")
	}

	log.Print("server shutting down...")
}
