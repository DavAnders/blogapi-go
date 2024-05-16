package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title string `bson:"title" json:"title" binding:"required"`
	Content string `bson:"content" json:"content" binding:"required"`
	PublishedAt time.Time `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"`
	AuthorID primitive.ObjectID `bson:"authorId" json:"authorId"`
	AuthorUsername string `bson:"authorUsername" json:"authorUsername"`
}

type PostResponse struct {
    ID             string    `json:"id"`
    Title          string    `json:"title"`
    Content        string    `json:"content"`
    PublishedAt    time.Time `json:"publishedAt"`
    AuthorID       string    `json:"authorId"`
    AuthorUsername string    `json:"authorUsername"`
}
