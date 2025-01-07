package main

import (
	"fmt"
	"go-server/src/handlers/categories"
	"go-server/src/handlers/products"
	"go-server/src/handlers/reviews"
	"log"
	"net/http"
)

// homeHandler handles requests to the root route "/".
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Go Server!")
}

// apiHandler handles requests to the "/api" route.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Go server api running successfully!")
}

// apiSubHandler handles requests to any sub-route under "/api/*".
func apiSubHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Fprintf(w, "You've reached the API sub-route: %s\n", path)
}

func main() {
	// Categories and Products Routes (These should be defined first)
	http.HandleFunc("/categories", categories.HandleCategories)
	http.HandleFunc("/categories/", categories.HandleCategory)
	http.HandleFunc("/products", products.HandleProducts)
	http.HandleFunc("/products/", products.HandleProduct)
	http.HandleFunc("/reviews", reviews.HandleReviews)
	http.HandleFunc("/reviews/", reviews.HandleReview)

	// General API Routes (Define after specific routes)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/api/", apiSubHandler)

	// Start the server
	port := "8080"
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
