package controller

import (
	"encoding/json"
	"net/http"

	"github.com/DavAnders/blogapi-go/pkg/jwt"
)

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := c.repo.ValidateCredentials(r.Context(), credentials.Username, credentials.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := jwt.GenerateToken(*user)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// func setTokenAsCookie(w http.ResponseWriter, tokenString string) {
//     http.SetCookie(w, &http.Cookie{
//         Name:     "token",
//         Value:    tokenString,
//         Expires:  time.Now().Add(1 * time.Hour),
//         HttpOnly: true, // JavaScript can't access cookie
//         Path:     "/",  // Cookie available on all paths
//         Secure:   false, // For development purposes
//         SameSite: http.SameSiteStrictMode,
//     })
// }
