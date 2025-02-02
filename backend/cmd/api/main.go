package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/routes"
	"github.com/mateuse/yahoo-fantasy-analyzer/internal/services"
)

func main() {

	// Load environment variables
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to MySQL
	err_sql := services.ConnectToMySQL()
	if err_sql != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err_sql)
	}
	// Create a new router
	router := mux.NewRouter()

	routes.RegisterHealthRoutes(router)
	routes.RegisterAuthRoutes(router)
	routes.RegisterYahooRoutes(router)
	routes.RegisterNHLRoutes(router)
	routes.RegisterStatsRoutes(router)
	routes.RegisterSearchRoutes(router)

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
