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

func initDB() {
	env := os.Getenv("ENV")
	if env == "" || env == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: No .env file found. Using system environment variables.")
		} else {
			log.Println(".env file loaded successfully.")
		}
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Error: DATABASE_URL environment variable is not set.")
	}

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

	router := http.NewServeMux()
	routes.RegisterRoutes(router, DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
