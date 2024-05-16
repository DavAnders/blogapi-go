package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DavAnders/blogapi-go/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id string) (*model.User, error)
	GetUsers(ctx context.Context) ([]UserProjection, error)
	ValidateCredentials(ctx context.Context, username, password string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
}

// UserProjection is a struct used to project only the necessary fields from a user
type UserProjection struct {
    Username  string    `bson:"username,omitempty" json:"username,omitempty"`
    CreatedAt time.Time `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
}

type userRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		db: db.Collection("users"),
	}
}

// Inserts a new user into the database
func (r *userRepository) CreateUser(ctx context.Context, user model.User) error {
	// Validate required fields
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	user.HashedPassword = string(hashedPassword)
	user.Password = "" // Clear the plain password
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	// Assign an ObjectID if not already specified
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	// Insert the user into the database
	result, err := r.db.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err
	}

	// Type assert the InsertedID and update the user's ID
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Failed to assert InsertedID to ObjectID")
		return fmt.Errorf("failed to assert InsertedID to ObjectID")
	}
	user.ID = oid

	log.Printf("Inserted user with ID: %v", user.ID.Hex())
	return nil
}


// Returns a single user from the database
func (r *userRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
    	return nil, fmt.Errorf("invalid ID format: %v", err)
	}
	err = r.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Returns all users from the database
func (r *userRepository) GetUsers(ctx context.Context) ([]UserProjection, error) {
	var users []UserProjection // Only project the necessary fields
	opts := options.Find().SetProjection(bson.M{"username": 1, "createdAt": 1})
	cur, err := r.db.Find(ctx, bson.M{}, opts)  // Can limit / sort if needed
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var user UserProjection
		if err := cur.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// ValidateCredentials checks a user's username and password against the stored values
func (r *userRepository) ValidateCredentials(ctx context.Context, username, password string) (*model.User, error) {
    log.Printf("Attempting to validate credentials for username: %s", username)

    var user model.User
    err := r.db.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            log.Printf("No user found for username: %s", username)
            return nil, fmt.Errorf("no user found with the given username")
        }
        log.Printf("Error retrieving user from database: %v", err)
        return nil, err
    }

    log.Printf("User found in database: %s, validating password", username)
    // Compare the stored hashed password with the provided password
    if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
        log.Printf("Password validation failed for user: %s", username)
        return nil, fmt.Errorf("invalid password")
    }

    log.Printf("Credentials validated successfully for user: %s", username)
    return &user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
    var user model.User
    err := r.db.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return model.User{}, fmt.Errorf("user not found")
        }
        return model.User{}, err
    }
    return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user model.User) error {
    user.UpdatedAt = time.Now()
    update := bson.M{
        "$set": bson.M{
            "bio":           user.Bio,
            "profilePicUrl": user.ProfilePicURL,
            "updatedAt":     user.UpdatedAt,
        },
    }

    filter := bson.M{"_id": user.ID}
    _, err := r.db.UpdateOne(ctx, filter, update)
    if err != nil {
        log.Printf("Error updating user: %v", err)
        return err
    }

    log.Printf("Updated user with ID: %v", user.ID.Hex())
    return nil
}