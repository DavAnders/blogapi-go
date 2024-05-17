package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepository struct {
    db *mongo.Collection
}

func NewAdminRepository(db *mongo.Database) *AdminRepository {
    return &AdminRepository{
        db: db.Collection("admins"),
    }
}

func (repo *AdminRepository) IsAdmin(ctx context.Context, userID string) (bool, error) {
    count, err := repo.db.CountDocuments(ctx, bson.M{"userId": userID})
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
