package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/middleware"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: could not load .env file, using system env variables")
	}

	postgresURL := os.Getenv("POSTGRES_URL")
	serverPort := os.Getenv("SERVER_PORT")
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	database, err := db.InitDB(postgresURL)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	handler := handlers.NewHandler(database)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Backend is running"))
	})
	mux.Handle("/v1/swift-codes/", middleware.RateLimitMiddleware(http.HandlerFunc(handler.SwiftHandler)))
	mux.Handle("/v1/swift-codes/country/", middleware.RateLimitMiddleware(http.HandlerFunc(handler.GetSwiftCodesByCountryHandler)))

	// cors configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	fmt.Printf("server listening on port %s\n", serverPort)
	if err := http.ListenAndServe(":"+serverPort, c.Handler(mux)); err != nil {
		log.Fatalf("server startup error: %v", err)
	}
}
