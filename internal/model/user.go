package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username" binding:"required"`
	Password string `bson:"-" json:"password"`
	HashedPassword string `bson:"hashedPassword" json:"-"`
	Email string `bson:"email" json:"email" binding:"required"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	Author bool `bson:"author" json:"author"`
	Bio            string             `bson:"bio,omitempty" json:"bio,omitempty"`
    ProfilePicURL  string             `bson:"profilePicUrl,omitempty" json:"profilePicUrl,omitempty"`
    UpdatedAt      time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}