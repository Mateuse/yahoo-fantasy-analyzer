package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/routes"
)

func main() {

	// Load environment variables
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	clientID := os.Getenv("YAHOO_CLIENT_ID")
	if clientID == "" {
		log.Fatal("CLIENT_ID environment variable is required")
	}
	// Create a new router
	router := mux.NewRouter()

	routes.RegisterHealthRoutes(router)
	routes.RegisterAuthRoutes(router)

	// Define allowed CORS options
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173", "https://your-frontend.com"}), // Update with your allowed origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),            // HTTP methods allowed
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),                      // Headers allowed
	)

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsOptions(router)))
}
