package middleware

import (
	"context"
	"net/http"

	"github.com/DavAnders/blogapi-go/internal/repository"
)

func AdminMiddleware(repo repository.AdminRepository) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userID, ok := r.Context().Value("userID").(string)
            if !ok {
                http.Error(w, "Unauthorized access", http.StatusUnauthorized)
                return
            }

            // Check if the user is an admin
            isAdmin, err := repo.IsAdmin(context.Background(), userID)
            if err != nil {
                http.Error(w, "Failed to verify admin status", http.StatusInternalServerError)
                return
            }
            if !isAdmin {
                http.Error(w, "Admin privileges required", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
