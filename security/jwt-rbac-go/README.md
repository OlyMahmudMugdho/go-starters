# 🔐 Building a Secure JWT RBAC Authentication System in Go with PostgreSQL

## 🧩 Introduction

In modern web applications, **authentication** and **authorization** are essential. In this guide, we’ll build a simple but secure **JWT-based Role-Based Access Control (RBAC)** system using **Go**, **PostgreSQL**, **bcrypt** for password hashing, and native `net/http` package — without using any heavy frameworks.

Features:

- ✅ register users and store their hashed passwords
- ✅ log in and generate JWT tokens
- ✅ protect routes using middleware
- ✅ restrict access based on user roles
- ✅ test it with `curl` and a shell script

---

## 🗂️ Project Structure

Here's how our project is organized:

```
.
├── auth
│   └── auth.go            # Functions for token creation and validation
├── database
│   └── database.go        # PostgreSQL database connection
├── handlers
│   └── handlers.go        # HTTP handler functions (register, login, etc.)
├── middlewares
│   └── middlewares.go     # Auth and role-based middlewares
├── models
│   └── models.go          # User model
├── routes
│   └── routes.go          # Route registration
├── main.go                # Entry point of the app
├── go.mod / go.sum        # Go modules
└── test_auth.sh           # Shell script to test everything
```

---

## 🧱 Breakdown


---

### 1️⃣ Connecting to PostgreSQL

`database/database.go`

```go
package database

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

// InitDB initializes a PostgreSQL database connection
func InitDB() {
	var err error
	DB, err = sql.Open("postgres", "user=youruser password=yourpass dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("Database connected.")
}
```

✅ Replace the connection string with your actual database credentials.

---

### 2️⃣ User Model

`models/models.go`

```go
package models

// User represents a user in the system
type User struct {
	ID       int
	Username string
	Password string
	Role     string
}
```

---

### 3️⃣ Authentication Logic

`auth/auth.go`

```go
package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key") // Keep this secret

// Claims struct for JWT payload
type Claims struct {
	Username string
	Role     string
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for a user
func GenerateToken(username, role string) (string, error) {
	claims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}
```

---

### 4️⃣ Middleware for Auth and Role Check

`middlewares/middlewares.go`

```go
package middlewares

import (
	"net/http"
	"strings"
	"context"
	"your_project/auth"

	"github.com/golang-jwt/jwt/v5"
)

// Key to store user info in context
type ContextKey string
const userKey ContextKey = "user"

// AuthMiddleware validates JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store user info in context
		ctx := context.WithValue(r.Context(), userKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userKey).(*auth.Claims)
		if user.Role != role {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

---

### 5️⃣ Handlers

`handlers/handlers.go`

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"your_project/database"
	"your_project/models"
	"your_project/auth"

	"golang.org/x/crypto/bcrypt"
)

// Register a new user with hashed password
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Hashing error", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)",
		user.Username, string(hashed), user.Role)
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User registered successfully"))
}

// Login and return a JWT token
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User
	json.NewDecoder(r.Body).Decode(&req)

	var stored models.User
	err := database.DB.QueryRow("SELECT password, role FROM users WHERE username=$1", req.Username).
		Scan(&stored.Password, &stored.Role)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(req.Username, stored.Role)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Public route
func PublicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("🌐 Public route: Anyone can access this"))
}

// User-only route
func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("🔐 Welcome user!"))
}

// Admin-only route
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("👑 Welcome admin!"))
}
```

---

### 6️⃣ Routes Setup

`routes/routes.go`

```go
package routes

import (
	"net/http"
	"your_project/handlers"
	"your_project/middlewares"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)

	mux.HandleFunc("/public", handlers.PublicHandler)

	mux.Handle("/user", middlewares.AuthMiddleware(http.HandlerFunc(handlers.UserHandler)))
	mux.Handle("/admin", middlewares.AuthMiddleware(
		middlewares.RoleMiddleware("admin", http.HandlerFunc(handlers.AdminHandler))),
	)
}
```

---

### 7️⃣ Main File

`main.go`

```go
package main

import (
	"net/http"
	"your_project/database"
	"your_project/routes"
)

func main() {
	database.InitDB()

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)

	http.ListenAndServe(":8080", mux)
}
```

---

## 🧪 Testing with Shell Script

Create a file called `test_auth.sh` and run it to automate testing:

```bash
chmod +x test_auth.sh
./test_auth.sh
```

It will:

- Register a user and an admin
- Log in and get tokens
- Test `/public`, `/user`, and `/admin` routes
- Display results

✅ Things you must do:  
**🌍 Visit [https://olymahmud.vercel.app](https://olymahmud.vercel.app)**
