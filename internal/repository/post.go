package repository

import (
	"context"
	"fmt"

	"github.com/DavAnders/blogapi-go/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Interface for querying posts from db
type PostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPosts(ctx context.Context, filter bson.M, limit int64, skip int64) ([]model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	UpdatePost(ctx context.Context, post model.Post) error
	DeletePost(ctx context.Context, id string, userID *string) error
    GetPostsByUser(ctx context.Context, userID string) ([]model.Post, error)
}

type postRepository struct {
	db *mongo.Collection
}

// Create a new post repository
func NewPostRepository(db *mongo.Database) PostRepository {
	return &postRepository{
		db: db.Collection("posts"),
	}
}

// Inserts a new post into the database
func (r *postRepository) CreatePost(ctx context.Context, post *model.Post) error {
	result, err := r.db.InsertOne(ctx, post)
    if err != nil {
        return err
    }
    post.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

// Returns posts from the database with optional filters, limit, and skip
func (r *postRepository) GetPosts(ctx context.Context, filter bson.M, limit int64, skip int64) ([]model.Post, error) {
    var posts []model.Post
    opts := options.Find().
        SetLimit(limit).
        SetSkip(skip).
        SetSort(bson.D{{Key: "publishedAt", Value: -1}}) // Use keyed fields

    cur, err := r.db.Find(ctx, filter, opts)
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)
    for cur.Next(ctx) {
        var post model.Post
        if err := cur.Decode(&post); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    if err := cur.Err(); err != nil {
        return nil, err
    }
    return posts, nil
}


// Find a post by its ID
func (r *postRepository) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
    var post model.Post
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err  // If the ID is not a valid ObjectId
    }
    if err := r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&post); err != nil {
        return nil, err
    }
    return &post, nil
}

// Updates a post in the database
func (r *postRepository) UpdatePost(ctx context.Context, post model.Post) error {
    update := bson.M{
        "$set": bson.M{
            "title":   post.Title,
            "content": post.Content,
        },
    }

    filter := bson.M{"_id": post.ID}
    result, err := r.db.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    if result.MatchedCount == 0 {
        return fmt.Errorf("no post found with given ID")
    }
    return nil
}

// Deletes a post from the database
func (r *postRepository) DeletePost(ctx context.Context, id string, userID *string) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err  // If the ID is not a valid ObjectId
    }
    filter := bson.M{"_id": objID}
    if userID != nil {
        filter["authorId"] = *userID  // Add author check only if userID is provided
    }

    result, err := r.db.DeleteOne(ctx, filter)
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        if userID != nil {
            return fmt.Errorf("no post found with given ID or unauthorized")
        }
        return fmt.Errorf("no post found with given ID")
    }
    return nil
}

// Get all posts by a user
func (r *postRepository) GetPostsByUser(ctx context.Context, userID string) ([]model.Post, error) {
    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return nil, err
    }

    filter := bson.M{"authorId": objID}
    options := options.Find().SetLimit(5) // Set limit to 5

    cursor, err := r.db.Find(ctx, filter, options)
    if err != nil {
        return nil, err
    }

    var posts []model.Post
    if err := cursor.All(ctx, &posts); err != nil {
        return nil, err
    }

    return posts, nil
}
