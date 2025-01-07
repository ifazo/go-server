package categories

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"go-server/src/models"
	"go-server/src/data"
)

// Handle Categories route
func HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetCategories(w)
	case http.MethodPost:
		CreateCategory(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle individual category route
func HandleCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetCategory(w, id)
	case http.MethodDelete:
		DeleteCategory(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get all categories
func GetCategories(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.Categories)
}

// Get a single category by ID
func GetCategory(w http.ResponseWriter, id int) {
	for _, category := range data.Categories {
		if category.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}

// Create a new category
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	category.ID = data.NextCategoryID
	data.NextCategoryID++
	category.Items = []models.Product{}
	data.Categories = append(data.Categories, category)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// Delete a category
func DeleteCategory(w http.ResponseWriter, id int) {
	for i, category := range data.Categories {
		if category.ID == id {
			data.Categories = append(data.Categories[:i], data.Categories[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}
