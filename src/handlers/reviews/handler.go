package reviews

import (
	"context"
	"encoding/json"
	"fmt"
	"go-server/src/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

// HandleReviews handles requests to the "/reviews" route.
func HandleReviews(w http.ResponseWriter, r *http.Request, DB *pgx.Conn) {
	switch r.Method {
	case http.MethodGet:
		GetReviews(w, DB)
	case http.MethodPost:
		CreateReview(w, r, DB)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleReview handles requests to "/reviews/{id}" route.
func HandleReview(w http.ResponseWriter, r *http.Request, DB *pgx.Conn) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/reviews/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		GetReview(w, DB, id)
	case http.MethodPatch:
		UpdateReview(w, r, DB, id)
	case http.MethodDelete:
		DeleteReview(w, DB, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetReviews retrieves all reviews from the database.
func GetReviews(w http.ResponseWriter, DB *pgx.Conn) {
	rows, err := DB.Query(context.Background(), "SELECT id, product_id, rating, comment FROM reviews")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching reviews: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(&review.ID, &review.ProductID, &review.Rating, &review.Comment); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning review: %v", err), http.StatusInternalServerError)
			return
		}
		reviews = append(reviews, review)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// GetReview retrieves a single review by ID from the database.
func GetReview(w http.ResponseWriter, DB *pgx.Conn, id int) {
	var review models.Review
	err := DB.QueryRow(context.Background(), "SELECT id, product_id, rating, comment FROM reviews WHERE id=$1", id).
		Scan(&review.ID, &review.ProductID, &review.Rating, &review.Comment)

	if err != nil {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// CreateReview creates a new review and associates it with a product.
func CreateReview(w http.ResponseWriter, r *http.Request, DB *pgx.Conn) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := DB.QueryRow(
		context.Background(),
		"INSERT INTO reviews (product_id, rating, comment) VALUES ($1, $2, $3) RETURNING id",
		review.ProductID, review.Rating, review.Comment,
	).Scan(&review.ID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating review: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

// UpdateReview updates a review by ID in the database.
func UpdateReview(w http.ResponseWriter, r *http.Request, DB *pgx.Conn, id int) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	commandTag, err := DB.Exec(
		context.Background(),
		"UPDATE reviews SET product_id=$1, rating=$2, comment=$3 WHERE id=$4",
		review.ProductID, review.Rating, review.Comment, id,
	)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating review: %v", err), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(review)
}

// DeleteReview deletes a review by ID from the database.
func DeleteReview(w http.ResponseWriter, DB *pgx.Conn, id int) {
	commandTag, err := DB.Exec(context.Background(), "DELETE FROM reviews WHERE id=$1", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting review: %v", err), http.StatusInternalServerError)
		return
	}

	if commandTag.RowsAffected() == 0 {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
