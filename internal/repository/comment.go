package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/DavAnders/blogapi-go/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment model.Comment) error
	GetCommentsByPost(ctx context.Context, postID primitive.ObjectID) ([]model.Comment, error)
	UpdateComment(ctx context.Context, id string, userID string, comment model.Comment) error
	DeleteComment(ctx context.Context, id string, userID *primitive.ObjectID) error
}

type commentRepository struct {
	db *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) CommentRepository {
	return &commentRepository{
		db: db.Collection("comments"),
	}
}

func (r *commentRepository) CreateComment(ctx context.Context, comment model.Comment) error {
    comment.ID = primitive.NewObjectID()
    comment.CreatedAt = time.Now()
    _, err := r.db.InsertOne(ctx, comment)
    return err
}


func (r *commentRepository) GetCommentsByPost(ctx context.Context, postID primitive.ObjectID) ([]model.Comment, error) {
	var comments []model.Comment
	filter := bson.M{"postId": postID}
	cur, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var comment model.Comment
		if err := cur.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) UpdateComment(ctx context.Context, id string, userID string, comment model.Comment) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err  // If the ID is not a valid ObjectId
    }

    update := bson.M{"$set": bson.M{
        "content": comment.Content,
        "email": comment.Email, // Updating email for now, but might want to change this
        "updatedAt": time.Now(),
    }}
    filter := bson.M{"_id": objID, "author": userID} // Ensure that the author matches the userID

    result, err := r.db.UpdateOne(ctx, filter, update)
    if err != nil {
        return err
    }
    if result.MatchedCount == 0 {
        return fmt.Errorf("no comment found with given ID or unauthorized")
    }
    return nil
}

func (r *commentRepository) DeleteComment(ctx context.Context, id string, userID *primitive.ObjectID) error {
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err  // If the ID is not a valid ObjectId
    }

    filter := bson.M{"_id": objID}
    if userID != nil {
        filter["authorId"] = *userID // Add author check only if userID is provided
    }

    result, err := r.db.DeleteOne(ctx, filter)
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        if userID != nil {
            return fmt.Errorf("no comment found with given ID or unauthorized")
        }
        return fmt.Errorf("no comment found with given ID")
    }
    return nil
}


