package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DavAnders/blogapi-go/internal/api/middleware"
	"github.com/DavAnders/blogapi-go/internal/model"
	"github.com/DavAnders/blogapi-go/internal/repository"
	"github.com/go-chi/chi/v5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostController struct {
	repo repository.PostRepository
}

func NewPostController(repo repository.PostRepository) *PostController {
	return &PostController{
		repo: repo,
	}
}

// Handles POST requests
func (c *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
    var post model.Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    post.PublishedAt = time.Now()

    userID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Unauthorized - User ID missing or invalid", http.StatusUnauthorized)
        return
    }

    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        http.Error(w, "Unauthorized - User ID conversion error: "+err.Error(), http.StatusUnauthorized)
        return
    }
    post.AuthorID = objID

    username, ok := r.Context().Value(middleware.UsernameKey).(string)
    if !ok {
        http.Error(w, "Unauthorized - Missing username", http.StatusUnauthorized)
        return
    }
    post.AuthorUsername = username

    if err := c.repo.CreatePost(r.Context(), &post); err != nil {
        http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
        return
    }

    post.AuthorUsername = username
    
    response := model.PostResponse{
        ID:             post.ID.Hex(),
        Title:          post.Title,
        Content:        post.Content,
        PublishedAt:    post.PublishedAt,
        AuthorID:       post.AuthorID.Hex(), 
        AuthorUsername: post.AuthorUsername,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Handles GET requests to retrieve all posts
func (c *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
    limit := int64(10) // Default limit
    skip := int64(0)   // Default skip
    limitQuery := r.URL.Query().Get("limit")
    skipQuery := r.URL.Query().Get("skip")
    if limitQuery != "" {
        limit, _ = strconv.ParseInt(limitQuery, 10, 64)
    }
    if skipQuery != "" {
        skip, _ = strconv.ParseInt(skipQuery, 10, 64)
    }

    posts, err := c.repo.GetPosts(context.Background(), bson.M{}, limit, skip)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}


// Handles GET requests to retrieve a post by ID
func (c *PostController) GetPostByID(w http.ResponseWriter, r *http.Request) {
    postID := chi.URLParam(r, "id")
    if postID == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    post, err := c.repo.GetPostByID(context.Background(), postID)
    if err != nil {
        http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

// Handles GET requests to retrieve all posts by a user
func (c *PostController) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "userID")
    if userID == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    authUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Check if the requested userID matches the authenticated user's ID
    if userID != authUserID {
        http.Error(w, "Unauthorized - You can only view your own posts", http.StatusUnauthorized)
        return
    }

    posts, err := c.repo.GetPostsByUser(context.Background(), userID)
    if err != nil {
        http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}



// Handles PUT requests to update a post
func (c *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
    postID := chi.URLParam(r, "id")
    if postID == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    var updatedPost model.Post
    if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    objID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        http.Error(w, "Invalid Post ID", http.StatusBadRequest)
        return
    }

    updatedPost.ID = objID

    if err := c.repo.UpdatePost(context.Background(), updatedPost); err != nil {
        http.Error(w, "Failed to update post", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedPost)
}

// Handles DELETE requests to delete a post
func (c *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
    postID := chi.URLParam(r, "id")
    userID, ok := r.Context().Value("userID").(string)
    if !ok || postID == "" {
        http.Error(w, "Unauthorized or bad request", http.StatusUnauthorized)
        return
    }

    // Pass userID for regular user deletes
    if err := c.repo.DeletePost(context.Background(), postID, &userID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}


// Deletes a post by ID, only if the user is an admin
func (c *PostController) AdminDeletePost(w http.ResponseWriter, r *http.Request) {
    postID := chi.URLParam(r, "id")
    if postID == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    // Pass nil as userID for admin deletes
    if err := c.repo.DeletePost(context.Background(), postID, nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
