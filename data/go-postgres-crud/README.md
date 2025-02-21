# Building an Enterprise-Grade Go PostgreSQL CRUD REST API

In this article, we’ll create a production-ready CRUD (Create, Read, Update, Delete) REST API using Go 1.23+ and PostgreSQL. We will use modern Go practices, a modular structure, database migrations with `golang-migrate`, input validation, and context handling. Let's dive into it.


## Project Overview

### What We’re Building

We’re creating a REST API to manage `users` with these endpoints:

- **POST /users**: Create a new user.
- **GET /users**: List all users.
- **GET /users/{id}**: Get a user by ID.
- **PUT /users/{id}**: Update a user.
- **DELETE /users/{id}**: Delete a user.

The API uses:

- **Go**: A fast, simple, modern language.
- **PostgreSQL**: A robust relational database.
- **http.NewServeMux**: Go’s built-in router.
- **golang-migrate**: For database schema changes.
- **validator**: For input validation.

### Why This Approach?

- **Modular Structure**: Keeps code organized.
- **Enterprise-Ready**: Includes migrations, validation, and error handling.
- **Scalable**: Easy to extend.

---

## Project Structure

```
go-postgres-crud/
├── cmd/
│   └── server/
│       └── main.go         # Entry point of the application
├── internal/
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── database/          # Database connection and migrations
│   │   ├── db.go
│   │   └── migrations/
│   │       ├── 000001_create_users_table.up.sql
│   │       └── 000001_create_users_table.down.sql
│   ├── handlers/          # HTTP request handlers
│   │   └── user.go
│   ├── models/            # Data models and validation
│   │   └── user.go
│   └── repository/        # Database operations
│       └── user.go
├── .env                   # Environment variables (optional)
├── go.mod                 # Go module file
└── README.md              # Project documentation
```

- **`cmd/server/`**: Where the application starts.
- **`internal/`**: Private code split into packages.
- **`migrations/`**: SQL files for database changes.

---

## Step-by-Step Code Explanation

### 1. Initialization (`go.mod`)
```bash
# Creates a new Go module
go mod init go-postgres-crud

# Installs libraries for validation, PostgreSQL, and migrations
go get github.com/go-playground/validator/v10
go get github.com/lib/pq
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
go get github.com/joho/godotenv
```

#### Explanation
This step sets up the project as a Go module named `go-postgres-crud`. The `go mod init` command creates a `go.mod` file to track dependencies. The `go get` commands download essential libraries: `validator` for input checking, `lib/pq` for PostgreSQL connectivity, and `golang-migrate` components for managing database changes. These ensure the project has the necessary tools.

---

### 2. Configuration (`internal/config/config.go`)
```go
package config

import (
	"log"         // For printing messages to the console
	"os"          // To access environment variables
	"github.com/joho/godotenv" // Loads .env file for easy config
)

type Config struct {
	// Fields to store database and server settings
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func LoadConfig() *Config {
	// Try to load a .env file; if it fails, use environment variables instead
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Create and return a Config object with values from env vars or defaults
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),     // Database host
		DBPort:     getEnv("DB_PORT", "5432"),          // Database port
		DBUser:     getEnv("DB_USER", "postgres"),      // Database username
		DBPassword: getEnv("DB_PASSWORD", "mysecretpassword"), // Database password
		DBName:     getEnv("DB_NAME", "go_crud"),       // Database name
		ServerPort: getEnv("SERVER_PORT", "8080"),      // API server port
	}
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	// Check if the env var exists; if yes, return it
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	// If not, return the default value
	return defaultValue
}
```

#### Explanation
This file manages configuration settings like database credentials and server port. The `Config` struct holds these values. `LoadConfig` attempts to load them from a `.env` file (useful for local development) using `godotenv`, falling back to environment variables if the file isn’t found. The `getEnv` function retrieves each setting (e.g., `DB_HOST`) from the environment or uses a default (e.g., `localhost`) if unset. This approach keeps the configuration flexible and secure.

Example `.env`:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=mysecretpassword
DB_NAME=go_crud
SERVER_PORT=8080
```

---

### 3. Database Setup (`internal/database/db.go`)
```go
package database

import (
	"database/sql"                           // Go’s SQL database tools
	"fmt"                                    // For formatting strings and errors
	"log"                                    // For logging messages
	"go-postgres-crud/internal/config"       // To use config settings
	_ "github.com/lib/pq"                    // PostgreSQL driver (blank import)
	"github.com/golang-migrate/migrate/v4"            // Migration library
	"github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL migration driver
	"github.com/golang-migrate/migrate/v4/source/file"       // Reads migration files
)

// Creates a new database connection
func NewDB(cfg *config.Config) (*sql.DB, error) {
	// Build the connection string using config values
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Open a connection to PostgreSQL
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		// If it fails, return an error with details
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection to ensure it works
	if err := db.Ping(); err != nil {
		// If ping fails, return an error
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Return the working database connection
	return db, nil
}

// Applies migration files to update the database
func RunMigrations(db *sql.DB) error {
	// Load migration files from the migrations folder
	source, err := (&file.File{}).Open("file://internal/database/migrations")
	if err != nil {
		// If the folder can’t be opened, return an error
		return fmt.Errorf("failed to open migration source: %w", err)
	}

	// Set up a PostgreSQL driver for migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		// If driver setup fails, return an error
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create a migration instance with files and database connection
	m, err := migrate.NewWithInstance("file", source, "postgres", driver)
	if err != nil {
		// If initialization fails, return an error
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	// Run all "up" migrations (e.g., create tables)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// If it fails (and it’s not just "no changes"), return an error
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Log success if migrations run or no changes were needed
	log.Println("Database migrations applied successfully")
	return nil // No error means success
}
```

#### Explanation
This file handles database connectivity and schema updates. `NewDB` connects to PostgreSQL by constructing a connection string from config values (e.g., `host=localhost port=5432 ...`), opening the connection with `sql.Open`, and verifying it with `db.Ping`. If it fails (e.g., wrong credentials), it returns an error. `RunMigrations` uses `golang-migrate` to apply SQL files from `internal/database/migrations/`, setting up a PostgreSQL driver and running "up" migrations (e.g., table creation). It logs success and handles cases where no new changes apply.

---

### 4. Migration Files (`internal/database/migrations/`)
- **`000001_create_users_table.up.sql`**:
```sql
-- Creates the users table with four columns
CREATE TABLE users (
    id SERIAL PRIMARY KEY,          -- Auto-incrementing ID as the primary key
    name VARCHAR(100) NOT NULL,     -- Name field, can’t be empty
    email VARCHAR(100) UNIQUE NOT NULL, -- Email field, unique and required
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp of creation
);
```

- **`000001_create_users_table.down.sql`**:
```sql
-- Deletes the users table if we need to undo the migration
DROP TABLE users;
```

#### Explanation
These SQL files define database changes. The `.up.sql` file creates a `users` table with an auto-incrementing `id`, a required `name` (up to 100 characters), a unique and required `email`, and a `created_at` timestamp defaulting to the current time. The `.down.sql` file reverses this by dropping the table. These migrations ensure the database structure matches the application’s needs.

---

### 5. Models (`internal/models/user.go`)
```go
package models

import (
	"time"                    // For the created_at timestamp
	"github.com/go-playground/validator/v10" // For input validation
)

type User struct {
	// Fields with JSON and database tags
	ID        int       `json:"id" db:"id"`              // User ID
	Name      string    `json:"name" db:"name" validate:"required,min=2,max=100"` // Name with validation rules
	Email     string    `json:"email" db:"email" validate:"required,email"`       // Email with validation
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Creation timestamp
}

// Validates the User struct based on its tags
func (u *User) Validate() error {
	validate := validator.New() // Create a new validator instance
	return validate.Struct(u)   // Check the struct; return error if invalid
}
```

#### Explanation
This file defines the `User` struct, representing a user in both the database and API responses. It includes `ID`, `Name`, `Email`, and `CreatedAt`, with tags for JSON serialization (`json:`), database mapping (`db:`), and validation (`validate:`). `Name` must be 2-100 characters and required, `Email` must be a valid email and required. The `Validate` method checks these rules using the `validator` library, returning an error if the data is invalid (e.g., missing email).

---

### 6. Repository (`internal/repository/user.go`)
```go
package repository

import (
	"context"         // For request cancellation/timeout
	"database/sql"    // For database operations
	"go-postgres-crud/internal/models" // User model
)

type UserRepository struct {
	db *sql.DB // Holds the database connection
}

// Creates a new repository instance
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at`
	// Execute query and scan returned values into the user struct
	return r.db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&user.ID, &user.CreatedAt)
}

// Fetches all users from the database
func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, email, created_at FROM users`
	rows, err := r.db.QueryContext(ctx, query) // Run the query
	if err != nil {
		return nil, err // Return error if query fails
	}
	defer rows.Close() // Ensure rows are closed after use
	var users []models.User // Slice to hold all users
	for rows.Next() {       // Loop through each row
		var u models.User
		// Scan row data into a User struct
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u) // Add user to the slice
	}
	return users, nil // Return the list of users
}

// Fetches a single user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	var u models.User
	// Execute query and scan result into the User struct
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err // Return error if no user or query fails
	}
	return &u, nil // Return the user
}

// Updates a user by ID
func (r *UserRepository) Update(ctx context.Context, id int, user *models.User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	// Execute the update query
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, id)
	if err != nil {
		return err // Return error if query fails
	}
	// Check if any rows were updated
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return sql.ErrNoRows // Return no rows error if nothing was updated
	}
	user.ID = id // Set the ID in the user struct
	return nil   // Success
}

// Deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	// Execute the delete query
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err // Return error if query fails
	}
	// Check if any rows were deleted
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return sql.ErrNoRows // Return no rows error if nothing was deleted
	}
	return nil // Success
}
```

#### Explanation
This file manages database operations for the `users` table. The `UserRepository` struct holds a database connection, and `NewUserRepository` initializes it. The methods handle CRUD actions: `Create` inserts a user and retrieves the `id` and `created_at`; `GetAll` fetches all users into a slice by looping through query results; `GetByID` retrieves one user by ID; `Update` modifies a user’s details; `Delete` removes a user. Each method uses `context.Context` for query control and checks for affected rows to handle missing records.

---

### 7. Handlers (`internal/handlers/user.go`)
```go
package handlers

import (
	"context"         // For request context
	"database/sql"    // For error handling (sql.ErrNoRows)
	"encoding/json"   // For JSON encoding/decoding
	"net/http"        // For HTTP handling
	"strconv"         // For converting strings to integers
	"go-postgres-crud/internal/models"     // User model
	"go-postgres-crud/internal/repository" // Repository
)

type UserHandler struct {
	repo *repository.UserRepository // Holds the repository
}

// Creates a new handler instance
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Handles POST /users to create a user
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Get request context
	var user models.User
	// Decode JSON request body into User struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validate the user data
	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Save the user to the database
	if err := h.repo.Create(ctx, &user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	// Set response headers and status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user) // Send the created user as JSON
}

// Handles GET /users to list all users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Get request context
	// Fetch all users from the database
	users, err := h.repo.GetAll(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	// Set response headers and send users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Handles GET /users/{id} to get one user
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Get request context
	// Parse the ID from the URL path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// Fetch the user by ID
	user, err := h.repo.GetByID(ctx, id)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}
	// Set response headers and send user as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Handles PUT /users/{id} to update a user
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Get request context
	// Parse the ID from the URL path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var user models.User
	// Decode JSON request body into User struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validate the user data
	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Update the user in the database
	if err := h.repo.Update(ctx, id, &user); err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	// Set response headers and send updated user as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Handles DELETE /users/{id} to delete a user
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Get request context
	// Parse the ID from the URL path
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// Delete the user from the database
	if err := h.repo.Delete(ctx, id); err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}
	// Set status to 204 (No Content) for successful deletion
	w.WriteHeader(http.StatusNoContent)
}
```

#### Explanation
This file defines HTTP handlers for the API’s endpoints. The `UserHandler` struct links to the repository, and `NewUserHandler` initializes it. Each method handles a request: `Create` decodes JSON, validates it, saves the user, and returns it; `GetAll` fetches and returns all users; `GetByID` parses an ID and returns one user; `Update` updates a user after validation; `Delete` removes a user. Errors are handled with HTTP status codes (e.g., 400 for bad input, 404 for not found) to provide clear feedback.

---

### 8. Main Application (`cmd/server/main.go`)
```go
package main

import (
	"context"         // For server context (not used here but imported for completeness)
	"log"             // For logging messages
	"net/http"        // For HTTP server and routing
	"go-postgres-crud/internal/config"     // Config package
	"go-postgres-crud/internal/database"   // Database package
	"go-postgres-crud/internal/handlers"   // Handlers package
	"go-postgres-crud/internal/repository" // Repository package
)

// Entry point for the application
func main() {
	// Load configuration settings
	cfg := config.LoadConfig()

	// Connect to the database
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err) // Exit if connection fails
	}
	defer db.Close() // Close the database connection when the app stops

	// Run database migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err) // Exit if migrations fail
	}

	// Create repository and handler instances
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// Set up HTTP router
	mux := http.NewServeMux()
	// Register routes with their handlers
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users", userHandler.GetAll)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	mux.HandleFunc("PUT /users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)

	// Configure and start the HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort, // Server address (e.g., :8080)
		Handler: mux,                  // Use the router
	}

	// Log startup message
	log.Printf("Server starting on port %s...", cfg.ServerPort)
	// Start the server and handle errors
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err) // Exit if server fails unexpectedly
	}
}
```

#### Explanation
This file is the application’s entry point. It loads the configuration, connects to the database, and runs migrations to ensure the schema is ready. It then creates a repository and handler, sets up an HTTP router with `http.NewServeMux`, and maps each endpoint to its handler function. The server starts on the configured port (e.g., 8080), logging its status. If errors occur during startup, it exits with a message; `defer db.Close()` ensures the database connection is cleaned up when the app stops.

---

## Running the Project

1. **Set Up PostgreSQL**:
   - Create the database: `CREATE DATABASE go_crud;`.
   - Update `.env` with your credentials.

2. **Run the App**:
   ```bash
   go run cmd/server/main.go
   ```
   - The server starts on `http://localhost:8080`.

3. **Test with HTTP Client**:
   Use an `api-call.http` file:
   ```http
   # Create a user
   POST http://localhost:8080/users
   Content-Type: application/json
   {"name": "John Doe", "email": "john.doe@example.com"}
   ```

---

## Conclusion

This project demonstrates how to build a robust REST API with Go and PostgreSQL. Key takeaways:
- **Modularity**: Code is split into packages for clarity and reuse.
- **Migrations**: `golang-migrate` keeps the database in sync.
- **Validation**: Ensures data integrity.
- **Scalability**: Easy to add new features.

Extend it with authentication, logging, or more models. Happy coding!
