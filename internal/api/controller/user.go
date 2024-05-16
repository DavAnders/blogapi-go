package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DavAnders/blogapi-go/internal/api/middleware"
	"github.com/DavAnders/blogapi-go/internal/model"
	"github.com/DavAnders/blogapi-go/internal/repository"
	"github.com/DavAnders/blogapi-go/pkg/jwt"
	"github.com/go-chi/chi/v5"
)

type UserController struct {
	repo repository.UserRepository
}

func NewUserController(repo repository.UserRepository) *UserController {
	return &UserController{
		repo: repo,
	}
}

// Create user with JWT token
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
    var user model.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Check if the username already exists
    _, err := c.repo.GetUserByUsername(r.Context(), user.Username)
    if err == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusConflict)
        json.NewEncoder(w).Encode(map[string]string{"error": "Username already exists"})
        return
    } else if err.Error() != "user not found" {
        http.Error(w, "Failed to check user existence", http.StatusInternalServerError)
        return
    }

    // Create the new user
    err = c.repo.CreateUser(r.Context(), user)
    if err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    // Retrieve the newly created user from the database
    createdUser, err := c.repo.GetUserByUsername(r.Context(), user.Username)
    if err != nil {
        http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
        return
    }

    // Generate the JWT token
    token, err := jwt.GenerateToken(createdUser)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    // Set the token as a cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    token,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        Secure:   true,
        Path:     "/",
        SameSite: http.SameSiteStrictMode,
    })

    // Return the token in the response body as well
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
    })
}

// Handles POST requests to create a new user
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user model.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, fmt.Sprintf("Invalid user data: %v", err), http.StatusBadRequest)
        return
    }

    if err := c.repo.CreateUser(context.Background(), user); err != nil {
        log.Printf("Failed to create user: %v", err)  
        http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

// Handles GET requests to retrieve a single user
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := c.repo.GetUser(context.Background(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Handles GET requests to retrieve all users
func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.repo.GetUsers(context.Background())
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (c *UserController) GetUserProfile(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok || userID == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    user, err := c.repo.GetUser(r.Context(), userID)
    if err != nil {
        log.Println("Failed to retrieve user:", err)
        http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}


func (c *UserController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok || userID == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var updatedFields struct {
        Bio           string `json:"bio"`
        ProfilePicURL string `json:"profilePicUrl"`
    }
    if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := c.repo.GetUser(r.Context(), userID)
    if err != nil {
        log.Println("Failed to retrieve user:", err)
        http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
        return
    }

    user.Bio = updatedFields.Bio
    user.ProfilePicURL = updatedFields.ProfilePicURL

    if err := c.repo.UpdateUser(r.Context(), *user); err != nil {
        http.Error(w, "Failed to update profile", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

