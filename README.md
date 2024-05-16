# Go Blog API

This Go Blog API is a simple backend service designed for handling blog-related operations, including user authentication, post creation, and comment management. Built with Go using the Chi router and MongoDB, this API provides a solid foundation for a blog platform, offering RESTful endpoints for managing blog posts, user profiles, and comments.

## Features

**User Authentication**: Register and manage user sessions.  
**Blog Posts**: Create, update, delete, and retrieve blog posts.  
**Comments**: Add and manage comments on blog posts.  
**Profiles**: User profile management, allowing for updating bio and profile picture URLs.  
**Admins**: Administrative endpoints for deleting posts and comments.

## Tech Stack

**Go**  
**Chi Router**: HTTP router for Go that supports RESTful routing.  
**MongoDB**: NoSQL database for storing data.  
**Docker**: Optional, for containerization and easy deployment.

## Getting Started

### Prerequisites

- Go (version 1.22.1 used here)
- A MongoDB URI
- An environment file (.env) with the following variables:
  - `MONGO_URI`: Your MongoDB connection string
  - `SECRET_KEY`: A secret key for signing JWTs

### Installation

1. Clone the repository:

   ```
   git clone https://github.com/DavAnders/blogapi-go
   ```

2. Navigate to the project directory:

   ```
   cd blogapi-go
   ```

3. Install the dependencies:

   ```
   go mod download
   ```

4. Create an environment file (.env) and set the required variables:

   ```
   MONGO_URI=<your_mongodb_uri>
   SECRET_KEY=<your_secret_key>
   ```

5. Build and run the application:

   ```
   go build ./cmd/api/main.go && ./main
   ```

6. The API should now be running on http://localhost:8080.

7. You can test the API using tools like cURL or Postman.

### Docker Deployment (Optional)

1. Build the Docker image:

   ```
   docker build -t blogapi-go .
   ```

2. Run the Docker container:

   ```
   docker run -p 8080:8080 -d blogapi-go
   ```

3. The API should now be running on http://localhost:8080.

4. You can test the API using tools like cURL or Postman.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please feel free to open a pull request.
