package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/YahiaJouini/careflow/api/routes"
	"github.com/YahiaJouini/careflow/internal/config"
	"github.com/YahiaJouini/careflow/internal/db"
	"github.com/rs/cors"
)

func main() {
	config.LoadEnv()
	port, portError := config.GetEnv("PORT")
	mode, modeError := config.GetEnv("MODE")
	if portError != nil || modeError != nil {
		log.Fatal(portError, modeError)
	}
	if mode == "development" {
		port = ":5000"
	}

	instance := db.InitializeDB()
	defer instance.Close()
	db.Migrate()

	router := routes.InitializeRoutes()
	fmt.Println("Server running on port", port)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:4173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	// apply cors middleware to the handler
	corsRouter := corsHandler.Handler(router)

	// to apply security middleware
	secureRouter := routes.UseSecurityMiddleware(corsRouter)
	log.Fatal(http.ListenAndServe(port, secureRouter))
}
