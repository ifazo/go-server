package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"go-server/src/routes"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

// initDB initializes the database connection.
func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")
	var errConn error
	DB, errConn = pgx.Connect(context.Background(), connStr)
	if errConn != nil {
		log.Fatalf("Unable to connect to database: %v\n", errConn)
	}
	log.Println("Database connected successfully!")
}

func main() {
	initDB()
	defer DB.Close(context.Background())

	// Initialize routes
	router := http.NewServeMux()
	routes.RegisterRoutes(router, DB)

	// Start the server
	port := "8080"
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
