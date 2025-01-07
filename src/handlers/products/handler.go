package products

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"go-server/src/models"
	"go-server/src/data"
)

// Handle Products route
func HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetProducts(w)
	case http.MethodPost:
		CreateProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle individual product route
func HandleProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		GetProduct(w, id)
	case http.MethodPut:
		UpdateProduct(w, r, id)
	case http.MethodDelete:
		DeleteProduct(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get all products
func GetProducts(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.Products)
}

// Get a single product by ID
func GetProduct(w http.ResponseWriter, id int) {
	for _, product := range data.Products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}

// Create a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	product.ID = data.NextProductID
	data.NextProductID++
	product.Reviews = []models.Review{}
	data.Products = append(data.Products, product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// Update an existing product
func UpdateProduct(w http.ResponseWriter, r *http.Request, id int) {
	var updatedProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for i, product := range data.Products {
		if product.ID == id {
			data.Products[i].Name = updatedProduct.Name
			data.Products[i].Price = updatedProduct.Price
			data.Products[i].Stock = updatedProduct.Stock
			data.Products[i].CategoryID = updatedProduct.CategoryID

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data.Products[i])
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

// Delete a product
func DeleteProduct(w http.ResponseWriter, id int) {
	for i, product := range data.Products {
		if product.ID == id {
			data.Products = append(data.Products[:i], data.Products[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Product not found", http.StatusNotFound)
}
