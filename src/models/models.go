package models

// Product represents the structure of a product
type Product struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID int       `json:"category_id"`
	Reviews    []Review  `json:"reviews"` // Reviews for this product
}

// Category represents the structure of a category
type Category struct {
	ID    int       `json:"id"`
	Name  string    `json:"name"`
	Items []Product `json:"items"` // Products under this category
}

// Review represents the structure of a review
type Review struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}
