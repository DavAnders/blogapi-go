package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Admin struct {
    UserID primitive.ObjectID `bson:"userId" json:"userId"`
}
