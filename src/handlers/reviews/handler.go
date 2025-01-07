package reviews

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"go-server/src/models"
	"go-server/src/data"
)

// HandleReviews handles requests to the "/reviews" route.
func HandleReviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetReviews(w)
	case http.MethodPost:
		CreateReview(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleReview handles requests to "/reviews/{id}" route.
func HandleReview(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/reviews/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetReview(w, id)
	case http.MethodDelete:
		DeleteReview(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetReviews retrieves all reviews.
func GetReviews(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.Reviews)
}

// GetReview retrieves a single review by ID.
func GetReview(w http.ResponseWriter, id int) {
	for _, review := range data.Reviews {
		if review.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(review)
			return
		}
	}
	http.Error(w, "Review not found", http.StatusNotFound)
}

// CreateReview creates a new review and associates it with a product.
func CreateReview(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	review.ID = data.NextReviewID
	data.NextReviewID++
	data.Reviews = append(data.Reviews, review)

	// Add review to the associated product.
	for i, product := range data.Products {
		if product.ID == review.ProductID {
			data.Products[i].Reviews = append(data.Products[i].Reviews, review)
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

// DeleteReview deletes a review by ID.
func DeleteReview(w http.ResponseWriter, id int) {
	for i, review := range data.Reviews {
		if review.ID == id {
			// Remove review from the reviews slice.
			data.Reviews = append(data.Reviews[:i], data.Reviews[i+1:]...)

			// Remove review from the associated product.
			for j, product := range data.Products {
				if product.ID == review.ProductID {
					for k, prodReview := range product.Reviews {
						if prodReview.ID == id {
							data.Products[j].Reviews = append(product.Reviews[:k], product.Reviews[k+1:]...)
							break
						}
					}
					break
				}
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Review not found", http.StatusNotFound)
}
