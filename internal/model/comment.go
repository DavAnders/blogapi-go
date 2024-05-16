package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID primitive.ObjectID `bson:"postId" json:"postId"`
	Author string `bson:"author" json:"author" binding:"required"`
	AuthorID primitive.ObjectID `bson:"authorId" json:"authorId"`
	Email string `bson:"email,omitempty" json:"email,omitempty"`
	Content string `bson:"content" json:"content" binding:"required"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}