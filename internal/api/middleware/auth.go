package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID string `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type ContextKey string // Used to store the user ID in the context of each request.

const (
    UserIDKey ContextKey = "userID"
    UsernameKey ContextKey = "username"
)

// AuthMiddleware validates the JWT token from the Authorization header and injects the user ID into the context.
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header is required", http.StatusUnauthorized)
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenStr == authHeader || tokenStr == "" {
            http.Error(w, "Invalid token format", http.StatusUnauthorized)
            return
        }

        token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            // Validate the alg is what we expect:
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(os.Getenv("SECRET_KEY")), nil
        })

        if err != nil {
            http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(*Claims)
        if !ok || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Inject user ID and username into the context of each request
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        ctx = context.WithValue(ctx, UsernameKey, claims.Username) 
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}