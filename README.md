
# Go Server Project

This is a simple Go server project that connects to a PostgreSQL database. It handles basic CRUD operations for products, categories, and reviews via RESTful APIs.

## Project Structure

```
go-server/
├── src/
│   ├── handlers/          
│   │   ├── categories.go  # Handlers for categories (GET, POST, DELETE)
│   │   ├── products.go    # Handlers for products (GET, POST, PUT, DELETE)
│   │   └── reviews.go     # Handlers for reviews (GET, POST, DELETE)
│   ├── models/
│   │   └── models.go      # Models for Product, Category, and Review
│   ├── routes/
│   │   └── routes.go      # Routes for Product, Category, and Review
│   └── main.go            # Main server entry point
├── go.mod                 # Go module file
├── go.sum                 # Go sum file
└── README.md              # Project documentation
```

## Requirements

- Go 1.18+ (or later)
- PostgreSQL database

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/ifazo/go-server.git
   cd go-server
   ```

2. **Set up PostgreSQL:**

   - Make sure you have PostgreSQL installed on your machine.
   - Create a new database and set up the tables (use SQL scripts in the project as examples).
   - Example SQL to create `categories` table:

     ```sql
     CREATE TABLE categories (
         id SERIAL PRIMARY KEY,
         name VARCHAR(100) NOT NULL
     );
     ```

3. **Set environment variables:**

   Create a `.env` file to store your database connection string.

   ```env
   DATABASE_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
   ```

4. **Install dependencies:**

   Run the following command to install Go dependencies:

   ```bash
   go mod tidy
   ```

## Running the Server

1. **Initialize the database connection** (done automatically when starting the server).
2. **Start the server** by running the following command:

   ```bash
   go run src/main.go
   ```

3. **The server will run on**: `http://localhost:8080`

## API Endpoints

- **GET /categories** - Retrieve all categories
- **GET /categories/{id}** - Retrieve a category by ID
- **POST /categories** - Create a new category
- **DELETE /categories/{id}** - Delete a category by ID

- **GET /products** - Retrieve all products
- **GET /products/{id}** - Retrieve a product by ID
- **POST /products** - Create a new product
- **PUT /products/{id}** - Update an existing product
- **DELETE /products/{id}** - Delete a product by ID

- **GET /reviews** - Retrieve all reviews
- **GET /reviews/{id}** - Retrieve a review by ID
- **POST /reviews** - Create a new review
- **DELETE /reviews/{id}** - Delete a review by ID

## Example Usage

### Fetch all Categories

```bash
curl -X GET http://localhost:8080/categories
```

### Create a new Product

```bash
curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name": "Product Name", "price": 99.99, "stock": 50, "category_id": 1}'
```

### Delete a Review

```bash
curl -X DELETE http://localhost:8080/reviews/1
```

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string for your local or remote PostgreSQL database.

Example:

```env
DATABASE_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
