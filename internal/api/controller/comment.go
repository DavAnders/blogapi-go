package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DavAnders/blogapi-go/internal/api/middleware"
	"github.com/DavAnders/blogapi-go/internal/model"
	"github.com/DavAnders/blogapi-go/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentController struct {
	repo repository.CommentRepository
}

func NewCommentController(repo repository.CommentRepository) *CommentController {
	return &CommentController{
		repo: repo,
	}
}

// Handles POST requests to create a new comment
func (c *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
    var comment model.Comment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        http.Error(w, "Invalid comment data", http.StatusBadRequest)
        return
    }

    // Get the user ID and username from the context
    userID, ok := r.Context().Value(middleware.UserIDKey).(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    username, ok := r.Context().Value(middleware.UsernameKey).(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusUnauthorized)
        return
    }

    postObjID, err := primitive.ObjectIDFromHex(comment.PostID.Hex())
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    comment.AuthorID = objID
    comment.Author = username
    comment.ID = primitive.NewObjectID()
    comment.CreatedAt = time.Now()
    comment.PostID = postObjID 

    if err := c.repo.CreateComment(context.Background(), comment); err != nil {
        http.Error(w, "Failed to create comment", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comment)
}


// Handles GET requests to retrieve all comments for a post
func (c *CommentController) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
    postID := chi.URLParam(r, "id")
    if postID == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    objID, err := primitive.ObjectIDFromHex(postID)
    if err != nil {
        http.Error(w, "Invalid Post ID", http.StatusBadRequest)
        return
    }

    comments, err := c.repo.GetCommentsByPost(context.Background(), objID)
    if err != nil {
        http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}


// Handles PUT requests to update a comment
func (c *CommentController) UpdateComment(w http.ResponseWriter, r *http.Request) {
    commentID := chi.URLParam(r, "id")
    if commentID == "" {
        http.Error(w, "Comment ID is required", http.StatusBadRequest)
        return
    }

    userID, ok := r.Context().Value("userID").(string)
    if !ok {
        http.Error(w, "Unauthorized or bad request", http.StatusUnauthorized)
        return
    }

    var comment model.Comment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    // Update the comment directly with user authorization check in the repo layer
    if err := c.repo.UpdateComment(context.Background(), commentID, userID, comment); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // Explicitly signify a successful update
    json.NewEncoder(w).Encode(comment)
}

// Handles DELETE requests to delete a comment
func (c *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
    commentID := chi.URLParam(r, "id")
    if commentID == "" {
        http.Error(w, "Comment ID is required", http.StatusBadRequest)
        return
    }

    // Extract userID from context
    userIDValue := r.Context().Value(middleware.UserIDKey)
    userID, ok := userIDValue.(string)
    if !ok || userID == "" {
        http.Error(w, "Unauthorized or bad request", http.StatusUnauthorized)
        return
    }

    // Convert userID to primitive.ObjectID
    objUserID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusUnauthorized)
        return
    }

    // Delete the comment directly with user authorization check in the repo layer
    if err := c.repo.DeleteComment(context.Background(), commentID, &objUserID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent) // No Content is typical for a successful delete operation
}



// Deletes a comment by ID, only if the user is an admin
func (c *CommentController) AdminDeleteComment(w http.ResponseWriter, r *http.Request) {
    commentID := chi.URLParam(r, "id")
    if commentID == "" {
        http.Error(w, "Comment ID is required", http.StatusBadRequest)
        return
    }

    // Pass nil as userID to indicate an admin deletion
    if err := c.repo.DeleteComment(context.Background(), commentID, nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

