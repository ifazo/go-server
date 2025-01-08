package routes

import (
	"go-server/src/handlers/categories"
	"go-server/src/handlers/products"
	"go-server/src/handlers/reviews"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// RegisterRoutes registers all the routes for the application.
func RegisterRoutes(router *http.ServeMux, db *pgx.Conn) {
	// Categories Routes
	router.HandleFunc("/api/categories", categories.HandleCategories(db))
	router.HandleFunc("/api/categories/", categories.HandleCategory(db))
	// Products Routes
	router.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		products.HandleProducts(w, r, db)
	})
	router.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		products.HandleProduct(w, r, db)
	})
	// Reviews Routes
	router.HandleFunc("/api/reviews", func(w http.ResponseWriter, r *http.Request) {
		reviews.HandleReviews(w, r, db)
	})
	router.HandleFunc("/api/reviews/", func(w http.ResponseWriter, r *http.Request) {
		reviews.HandleReview(w, r, db)
	})

	// API Route
	router.HandleFunc("/api", apiHandler)

	// General Route
	router.HandleFunc("/", homeHandler) // Must be registered last
}

// homeHandler handles requests to the root route "/".
func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("HomeHandler: %s", r.URL.Path)
	w.Write([]byte("Welcome to the Go Server!"))
}

// apiHandler handles requests to the "/api" route.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("APIHandler: %s", r.URL.Path)
	w.Write([]byte("Go server API running successfully!"))
}
