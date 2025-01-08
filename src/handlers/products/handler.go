package products

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	
	"go-server/src/models"

	"github.com/jackc/pgx/v5"
)

// HandleProducts handles requests to the "/products" route.
func HandleProducts(w http.ResponseWriter, r *http.Request, DB *pgx.Conn) {
	switch r.Method {
	case http.MethodGet:
		GetProducts(w, DB)
	case http.MethodPost:
		CreateProduct(w, r, DB)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleProduct handles requests to "/products/{id}" route.
func HandleProduct(w http.ResponseWriter, r *http.Request, DB *pgx.Conn) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		GetProduct(w, DB, id)
	case http.MethodPatch:
		UpdateProduct(w, r, DB, id)
	case http.MethodDelete:
		DeleteProduct(w, DB, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetProducts retrieves all products from the database.
func GetProducts(w http.ResponseWriter, DB *pgx.Conn) {
	rows, err := DB.Query(context.Background(), "SELECT id, name, price, stock, category_id FROM products")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching products: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning products: %v", err), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}

// GetProduct retrieves a product by its ID from the database.
func GetProduct(w http.ResponseWriter, DB *pgx.Conn, id int) {
	row := DB.QueryRow(context.Background(), "SELECT id, name, price, stock, category_id FROM products WHERE id = $1", id)

	var product models.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID); err != nil {
		http.Error(w, fmt.Sprintf("Error scanning product: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// CreateProduct creates a new product in the database.
func CreateProduct(w http.ResponseWriter, r *http.Request, DB *pgx.Conn) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding product: %v", err), http.StatusBadRequest)
		return
	}

	_, err := DB.Exec(context.Background(), "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4)", product.Name, product.Price, product.Stock, product.CategoryID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateProduct updates a product in the database.
func UpdateProduct(w http.ResponseWriter, r *http.Request, DB *pgx.Conn, id int) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding product: %v", err), http.StatusBadRequest)
		return
	}

	_, err := DB.Exec(context.Background(), "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5", product.Name, product.Price, product.Stock, product.CategoryID, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct deletes a product from the database.
func DeleteProduct(w http.ResponseWriter, DB *pgx.Conn, id int) {
	_, err := DB.Exec(context.Background(), "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting product: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}