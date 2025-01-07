package categories

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

// HandleCategories handles requests to the "/categories" route.
func HandleCategories(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetCategories(db, w)
		case http.MethodPost:
			CreateCategory(db, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// HandleCategory handles requests to "/categories/{id}" route.
func HandleCategory(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the ID from the path
		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
		idStr = strings.TrimSuffix(idStr, "/") // Remove trailing slash
		if idStr == "" {
			http.Error(w, "Category ID is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		// Handle request by method
		switch r.Method {
		case http.MethodGet:
			GetCategory(db, w, id)
		case http.MethodPatch:
			UpdateCategory(db, w, r, id)
		case http.MethodDelete:
			DeleteCategory(db, w, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// GetCategories retrieves all categories from the database.
func GetCategories(db *pgx.Conn, w http.ResponseWriter) {
	rows, err := db.Query(context.Background(), "SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			http.Error(w, "Failed to parse category data", http.StatusInternalServerError)
			return
		}
		categories = append(categories, map[string]interface{}{"id": id, "name": name})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GetCategory retrieves a single category by ID from the database.
func GetCategory(db *pgx.Conn, w http.ResponseWriter, id int) {
	var name string
	err := db.QueryRow(context.Background(), "SELECT name FROM categories WHERE id=$1", id).Scan(&name)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch category", http.StatusInternalServerError)
		}
		return
	}

	category := map[string]interface{}{"id": id, "name": name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// CreateCategory creates a new category in the database.
func CreateCategory(db *pgx.Conn, w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow(context.Background(), "INSERT INTO categories (name) VALUES ($1) RETURNING id", input.Name).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	category := map[string]interface{}{"id": id, "name": input.Name}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory updates a category by ID in the database.
func UpdateCategory(db *pgx.Conn, w http.ResponseWriter, r *http.Request, id int) {
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(context.Background(), "UPDATE categories SET name=$1 WHERE id=$2", input.Name, id)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	category := map[string]interface{}{"id": id, "name": input.Name}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory deletes a category by ID from the database.
func DeleteCategory(db *pgx.Conn, w http.ResponseWriter, id int) {
	result, err := db.Exec(context.Background(), "DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
