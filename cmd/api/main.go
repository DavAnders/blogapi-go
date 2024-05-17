package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/DavAnders/blogapi-go/internal/api/controller"
	"github.com/DavAnders/blogapi-go/internal/api/middleware"
	"github.com/DavAnders/blogapi-go/internal/repository"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file:", err)
    }

    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        log.Fatal("MONGO_URI is not set in .env file")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    if err = client.Ping(ctx, nil); err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }
    log.Println("Connected to MongoDB")

    postRepo := repository.NewPostRepository(client.Database("blogprod"))
    userRepo := repository.NewUserRepository(client.Database("blogprod"))
    commentRepo := repository.NewCommentRepository(client.Database("blogprod"))

    postController := controller.NewPostController(postRepo)
    userController := controller.NewUserController(userRepo)
    commentController := controller.NewCommentController(commentRepo)

    r := chi.NewRouter()
    r.Use(middleware.EnableCORS) // Global middleware

    // Public routes
    r.Group(func(r chi.Router) {
        r.Post("/login", userController.Login)
        r.Post("/register", userController.Register)
    })

    // API routes with AuthMiddleware
    r.Route("/api", func(r chi.Router) {
        r.Use(middleware.AuthMiddleware) // Applies to all routes in this route grouping

        r.Get("/posts", postController.GetPosts)
        r.Post("/posts", postController.CreatePost)
        r.Get("/posts/user/{userID}", postController.GetPostsByUser)
        r.Get("/posts/{id}", postController.GetPostByID)
        r.Put("/posts/{id}", postController.UpdatePost)
        r.Delete("/posts/{id}", postController.DeletePost)
        r.Get("/profile", userController.GetUserProfile)
        r.Put("/profile", userController.UpdateUserProfile)
        r.Get("/users", userController.GetUsers)
        r.Post("/users", userController.CreateUser)
        r.Get("/users/{id}", userController.GetUser)
        r.Post("/comments", commentController.CreateComment)
        r.Get("/comments/{id}", commentController.GetCommentsByPost)
        r.Put("/comments/{id}", commentController.UpdateComment)
        r.Delete("/comments/{id}", commentController.DeleteComment)
    })

    // Admin-specific routes
    r.Route("/admin", func(r chi.Router) {
        r.Use(middleware.AuthMiddleware) // General auth middleware
        r.Use(middleware.AdminMiddleware(*repository.NewAdminRepository(client.Database("blog")))) // Admin-specific middleware

        r.Delete("/posts/{id}", postController.AdminDeletePost)
        r.Delete("/comments/{id}", commentController.AdminDeleteComment)
    })

    // Starts the server listening on port 8080. In a production environment, consider:
    // - Using HTTPS for secure communication (using http.ListenAndServeTLS instead of http.ListenAndServe and obtaining a SSL/TLS certificate)
    // - Not listening on port 8080, which is commonly used for development. The default port for HTTPS is 443.
    // - Setting the port through an environment variable for flexibility.
    // - Using a reverse proxy like Nginx to handle HTTPS termination, load balancing, and other tasks.
    log.Println("Starting server on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}