package jwt

import (
	"log"
	"os"
	"time"

	"github.com/DavAnders/blogapi-go/internal/model"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
    Username string `json:"username"`
	UserID string `json:"userId"`
    jwt.StandardClaims
}

func getJWTKey() []byte {
    secretKey := os.Getenv("SECRET_KEY")
    if secretKey == "" {
        log.Fatal("SECRET_KEY is not set or is empty")
    }
    return []byte(secretKey)
}

// Generates a new JWT token
func GenerateToken(user model.User) (string, error) {
    jwtKey := getJWTKey()
    expirationTime := time.Now().Add(1 * time.Hour) // Token expires after 1 hour
    claims := &Claims{
        Username: user.Username,
		UserID: user.ID.Hex(), // Convert ObjectID to string
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    return tokenString, err
}

// Validates a JWT token
func ValidateToken(tokenString string) (*jwt.Token, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return getJWTKey(), nil
    })

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        log.Printf("User %s logged in", claims.Username)
    } else {
        log.Printf("Invalid token: %v", err)
        return nil, err
    }

    return token, nil
}


