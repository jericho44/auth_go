# Adding New Features Guide

This guide explains how to add new features to the JWT Authentication API following the established clean architecture patterns.

## 🏗️ Architecture Overview

The project follows clean architecture with these layers:

- **Models** - Data structures and domain entities
- **Repositories** - Data access layer (database operations)
- **Services** - Business logic and application rules
- **Controllers** - Request validation and response formatting
- **Handlers** - HTTP request/response handling
- **Routes** - URL routing and middleware setup

## 📋 Step-by-Step Process

### 1. Define the Model

Create or update models in the `models/` directory.

**Example: Adding a "Posts" feature**

```go
// models/post.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Title     string         `json:"title" gorm:"not null;size:200"`
    Content   string         `json:"content" gorm:"type:text"`
    UserID    uint           `json:"user_id" gorm:"not null;index"`
    User      User           `json:"user" gorm:"foreignKey:UserID"`
    Published bool           `json:"published" gorm:"default:false"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request/Response DTOs
type CreatePostRequest struct {
    Title   string `json:"title" validate:"required,min=1,max=200"`
    Content string `json:"content" validate:"required"`
}

type UpdatePostRequest struct {
    Title     *string `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
    Content   *string `json:"content,omitempty"`
    Published *bool   `json:"published,omitempty"`
}

type PostResponse struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Published bool      `json:"published"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. Create Repository Interface and Implementation

Define data access methods in the `repositories/` directory.

```go
// repositories/post_repository.go
package repositories

import (
    "auth-jwt/models"
    "gorm.io/gorm"
)

type PostRepositoryInterface interface {
    Create(post *models.Post) error
    GetByID(id uint) (*models.Post, error)
    GetByUserID(userID uint, limit, offset int) ([]*models.Post, error)
    Update(post *models.Post) error
    Delete(id uint) error
    GetPublished(limit, offset int) ([]*models.Post, error)
}

type PostRepository struct {
    db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepositoryInterface {
    return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
    return r.db.Create(post).Error
}

func (r *PostRepository) GetByID(id uint) (*models.Post, error) {
    var post models.Post
    err := r.db.Preload("User").First(&post, id).Error
    if err != nil {
        return nil, err
    }
    return &post, nil
}

func (r *PostRepository) GetByUserID(userID uint, limit, offset int) ([]*models.Post, error) {
    var posts []*models.Post
    err := r.db.Where("user_id = ?", userID).
        Limit(limit).Offset(offset).
        Order("created_at DESC").
        Find(&posts).Error
    return posts, err
}

func (r *PostRepository) Update(post *models.Post) error {
    return r.db.Save(post).Error
}

func (r *PostRepository) Delete(id uint) error {
    return r.db.Delete(&models.Post{}, id).Error
}

func (r *PostRepository) GetPublished(limit, offset int) ([]*models.Post, error) {
    var posts []*models.Post
    err := r.db.Where("published = ?", true).
        Preload("User").
        Limit(limit).Offset(offset).
        Order("created_at DESC").
        Find(&posts).Error
    return posts, err
}
```

### 3. Create Service Layer

Implement business logic in the `services/` directory.

```go
// services/post_service.go
package services

import (
    "errors"
    "auth-jwt/models"
    "auth-jwt/repositories"
)

type PostServiceInterface interface {
    CreatePost(userID uint, req models.CreatePostRequest) (*models.Post, error)
    GetPost(id uint) (*models.PostResponse, error)
    GetUserPosts(userID uint, page, limit int) ([]*models.PostResponse, error)
    UpdatePost(id, userID uint, req models.UpdatePostRequest) (*models.PostResponse, error)
    DeletePost(id, userID uint) error
    GetPublishedPosts(page, limit int) ([]*models.PostResponse, error)
}

type PostService struct {
    postRepo repositories.PostRepositoryInterface
    userRepo repositories.UserRepositoryInterface
}

func NewPostService(postRepo repositories.PostRepositoryInterface, userRepo repositories.UserRepositoryInterface) PostServiceInterface {
    return &PostService{
        postRepo: postRepo,
        userRepo: userRepo,
    }
}

func (s *PostService) CreatePost(userID uint, req models.CreatePostRequest) (*models.Post, error) {
    // Verify user exists
    user, err := s.userRepo.GetByID(userID)
    if err != nil || user == nil {
        return nil, errors.New("user not found")
    }

    post := &models.Post{
        Title:   req.Title,
        Content: req.Content,
        UserID:  userID,
    }

    if err := s.postRepo.Create(post); err != nil {
        return nil, errors.New("failed to create post")
    }

    return post, nil
}

func (s *PostService) GetPost(id uint) (*models.PostResponse, error) {
    post, err := s.postRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("post not found")
    }

    return &models.PostResponse{
        ID:        post.ID,
        Title:     post.Title,
        Content:   post.Content,
        Published: post.Published,
        Author:    post.User.Username,
        CreatedAt: post.CreatedAt,
        UpdatedAt: post.UpdatedAt,
    }, nil
}

func (s *PostService) UpdatePost(id, userID uint, req models.UpdatePostRequest) (*models.PostResponse, error) {
    post, err := s.postRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("post not found")
    }

    // Check ownership
    if post.UserID != userID {
        return nil, errors.New("unauthorized")
    }

    // Update fields if provided
    if req.Title != nil {
        post.Title = *req.Title
    }
    if req.Content != nil {
        post.Content = *req.Content
    }
    if req.Published != nil {
        post.Published = *req.Published
    }

    if err := s.postRepo.Update(post); err != nil {
        return nil, errors.New("failed to update post")
    }

    return &models.PostResponse{
        ID:        post.ID,
        Title:     post.Title,
        Content:   post.Content,
        Published: post.Published,
        Author:    post.User.Username,
        CreatedAt: post.CreatedAt,
        UpdatedAt: post.UpdatedAt,
    }, nil
}

func (s *PostService) DeletePost(id, userID uint) error {
    post, err := s.postRepo.GetByID(id)
    if err != nil {
        return errors.New("post not found")
    }

    // Check ownership
    if post.UserID != userID {
        return errors.New("unauthorized")
    }

    return s.postRepo.Delete(id)
}

// Additional methods...
```

### 4. Create Controller

Handle request validation and response formatting in `controllers/`.

```go
// controllers/post_controller.go
package controllers

import (
    "errors"
    "auth-jwt/models"
    "auth-jwt/services"
    "auth-jwt/utils"
)

type PostController struct {
    postService services.PostServiceInterface
}

type PostControllerResponse struct {
    Message string                 `json:"message,omitempty"`
    Post    *models.PostResponse   `json:"post,omitempty"`
    Posts   []*models.PostResponse `json:"posts,omitempty"`
}

func NewPostController(postService services.PostServiceInterface) *PostController {
    return &PostController{
        postService: postService,
    }
}

func (pc *PostController) CreatePost(userID int, req models.CreatePostRequest) (*PostControllerResponse, error) {
    // Validate request
    if err := utils.ValidateStruct(req); err != nil {
        return nil, errors.New("validation failed: " + err.Error())
    }

    post, err := pc.postService.CreatePost(uint(userID), req)
    if err != nil {
        return nil, err
    }

    return &PostControllerResponse{
        Message: "Post created successfully",
        Post: &models.PostResponse{
            ID:        post.ID,
            Title:     post.Title,
            Content:   post.Content,
            Published: post.Published,
            CreatedAt: post.CreatedAt,
            UpdatedAt: post.UpdatedAt,
        },
    }, nil
}

func (pc *PostController) GetPost(id int) (*PostControllerResponse, error) {
    post, err := pc.postService.GetPost(uint(id))
    if err != nil {
        return nil, err
    }

    return &PostControllerResponse{
        Post: post,
    }, nil
}

func (pc *PostController) UpdatePost(id, userID int, req models.UpdatePostRequest) (*PostControllerResponse, error) {
    // Validate request
    if err := utils.ValidateStruct(req); err != nil {
        return nil, errors.New("validation failed: " + err.Error())
    }

    post, err := pc.postService.UpdatePost(uint(id), uint(userID), req)
    if err != nil {
        return nil, err
    }

    return &PostControllerResponse{
        Message: "Post updated successfully",
        Post:    post,
    }, nil
}

func (pc *PostController) DeletePost(id, userID int) (*PostControllerResponse, error) {
    err := pc.postService.DeletePost(uint(id), uint(userID))
    if err != nil {
        return nil, err
    }

    return &PostControllerResponse{
        Message: "Post deleted successfully",
    }, nil
}
```

### 5. Create HTTP Handlers

Handle HTTP requests and responses in `handlers/`.

```go
// handlers/post.go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "auth-jwt/controllers"
    "auth-jwt/models"
    "auth-jwt/utils"

    "github.com/gorilla/mux"
)

var postController *controllers.PostController

func SetPostController(controller *controllers.PostController) {
    postController = controller
}

// @Summary Create a new post
// @Description Create a new post for the authenticated user
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.CreatePostRequest true "Post data"
// @Success 201 {object} controllers.PostControllerResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/posts [post]
func CreatePost(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value("user").(*utils.Claims)

    var req models.CreatePostRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.BadRequest(w, "Invalid request body", "Failed to parse JSON request")
        return
    }

    response, err := postController.CreatePost(claims.UserID, req)
    if err != nil {
        switch err.Error() {
        case "user not found":
            utils.NotFound(w, "User not found")
        default:
            if contains(err.Error(), "validation failed") {
                utils.ValidationError(w, err.Error(), "Please check your input")
            } else {
                utils.InternalServerError(w, "Failed to create post")
            }
        }
        return
    }

    utils.Created(w, response.Message, response.Post)
}

// @Summary Get a post by ID
// @Description Get a specific post by its ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} controllers.PostControllerResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/posts/{id} [get]
func GetPost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        utils.BadRequest(w, "Invalid post ID", "Post ID must be a number")
        return
    }

    response, err := postController.GetPost(id)
    if err != nil {
        switch err.Error() {
        case "post not found":
            utils.NotFound(w, "Post not found")
        default:
            utils.InternalServerError(w, "Failed to retrieve post")
        }
        return
    }

    utils.Success(w, "Post retrieved successfully", response.Post)
}

// @Summary Update a post
// @Description Update a post owned by the authenticated user
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body models.UpdatePostRequest true "Updated post data"
// @Success 200 {object} controllers.PostControllerResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/posts/{id} [put]
func UpdatePost(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value("user").(*utils.Claims)
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        utils.BadRequest(w, "Invalid post ID", "Post ID must be a number")
        return
    }

    var req models.UpdatePostRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.BadRequest(w, "Invalid request body", "Failed to parse JSON request")
        return
    }

    response, err := postController.UpdatePost(id, claims.UserID, req)
    if err != nil {
        switch err.Error() {
        case "post not found":
            utils.NotFound(w, "Post not found")
        case "unauthorized":
            utils.Forbidden(w, "You can only update your own posts")
        default:
            if contains(err.Error(), "validation failed") {
                utils.ValidationError(w, err.Error(), "Please check your input")
            } else {
                utils.InternalServerError(w, "Failed to update post")
            }
        }
        return
    }

    utils.Success(w, response.Message, response.Post)
}

// @Summary Delete a post
// @Description Delete a post owned by the authenticated user
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} controllers.PostControllerResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Security BearerAuth
// @Router /api/posts/{id} [delete]
func DeletePost(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value("user").(*utils.Claims)
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        utils.BadRequest(w, "Invalid post ID", "Post ID must be a number")
        return
    }

    response, err := postController.DeletePost(id, claims.UserID)
    if err != nil {
        switch err.Error() {
        case "post not found":
            utils.NotFound(w, "Post not found")
        case "unauthorized":
            utils.Forbidden(w, "You can only delete your own posts")
        default:
            utils.InternalServerError(w, "Failed to delete post")
        }
        return
    }

    utils.Success(w, response.Message, nil)
}

// Helper function
func contains(s, substr string) bool {
    return len(s) >= len(substr) && s[:len(substr)] == substr
}
```

### 6. Create Routes

Define URL routes in `routes/`.

```go
// routes/post.go
package routes

import (
    "auth-jwt/handlers"
    "auth-jwt/middleware"

    "github.com/gorilla/mux"
)

func SetupPostRoutes(r *mux.Router) {
    // Public routes
    r.HandleFunc("/posts/{id:[0-9]+}", handlers.GetPost).Methods("GET")
    r.HandleFunc("/posts", handlers.GetPublishedPosts).Methods("GET")

    // Protected routes
    protected := r.PathPrefix("/posts").Subrouter()
    protected.Use(middleware.AuthMiddleware)

    protected.HandleFunc("", handlers.CreatePost).Methods("POST")
    protected.HandleFunc("/{id:[0-9]+}", handlers.UpdatePost).Methods("PUT")
    protected.HandleFunc("/{id:[0-9]+}", handlers.DeletePost).Methods("DELETE")
    protected.HandleFunc("/my", handlers.GetMyPosts).Methods("GET")
}
```

### 7. Update Main Routes

Add your new routes to the main router.

```go
// routes/routes.go - Add to SetupRoutes function
func SetupRoutes() *mux.Router {
    r := mux.NewRouter()

    // Existing routes...
    SetupAuthRoutes(r.PathPrefix("/auth").Subrouter())
    SetupUserRoutes(r.PathPrefix("/api/user").Subrouter())

    // Add new routes
    SetupPostRoutes(r.PathPrefix("/api/posts").Subrouter())

    return r
}
```

### 8. Update Main.go

Initialize your new components in main.go.

```go
// main.go - Add to main function
func main() {
    // ... existing code ...

    // Initialize repositories
    userRepo := repositories.NewUserRepository(database.DB)
    authRepo := repositories.NewAuthRepository(database.DB)
    postRepo := repositories.NewPostRepository(database.DB) // Add this

    // Initialize services
    userService := services.NewUserService(userRepo)
    authService := services.NewAuthService(userRepo, authRepo)
    postService := services.NewPostService(postRepo, userRepo) // Add this

    // Initialize controllers
    userController := controllers.NewUserController(userService)
    authController := controllers.NewAuthController(authService)
    postController := controllers.NewPostController(postService) // Add this

    // Set controllers in handlers
    handlers.SetUserController(userController)
    handlers.SetAuthController(authController)
    handlers.SetPostController(postController) // Add this

    // ... rest of the code ...
}
```

### 9. Update Database Migration

Add your new model to the auto-migration.

```go
// database/db.go - Update AutoMigrate function
func AutoMigrate() error {
    err := DB.AutoMigrate(
        &models.User{},
        &models.PasswordResetToken{},
        &models.LoginAttempt{},
        &models.Post{}, // Add this
    )

    if err != nil {
        return fmt.Errorf("failed to auto-migrate: %w", err)
    }

    return nil
}
```

## 🧪 Testing Your New Feature

### 1. Test with cURL

```bash
# Create a post
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"My First Post","content":"This is the content of my post"}'

# Get a post
curl -X GET http://localhost:8080/api/posts/1

# Update a post
curl -X PUT http://localhost:8080/api/posts/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"Updated Title","published":true}'

# Delete a post
curl -X DELETE http://localhost:8080/api/posts/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 2. Update Swagger Documentation

Run swagger generation to update API docs:

```bash
swag init
```

## 📝 Best Practices

### 1. Error Handling

- Use consistent error messages
- Return appropriate HTTP status codes
- Log errors for debugging

### 2. Validation

- Validate all input data
- Use struct tags for validation rules
- Return clear validation error messages

### 3. Security

- Always validate user permissions
- Use middleware for authentication
- Sanitize input data

### 4. Database

- Use transactions for complex operations
- Add proper indexes for performance
- Use soft deletes when appropriate

### 5. Testing

- Write unit tests for services
- Test error scenarios
- Use integration tests for handlers

### 6. Documentation

- Add Swagger comments to handlers
- Update API documentation
- Document business logic in services

## 🔧 Common Patterns

### Pagination

```go
type PaginationRequest struct {
    Page  int `json:"page" validate:"min=1"`
    Limit int `json:"limit" validate:"min=1,max=100"`
}

func (s *PostService) GetPaginatedPosts(req PaginationRequest) ([]*models.PostResponse, *PaginationMeta, error) {
    offset := (req.Page - 1) * req.Limit
    posts, err := s.postRepo.GetPublished(req.Limit, offset)
    // ... implementation
}
```

### Search and Filtering

```go
type PostFilter struct {
    UserID    *uint   `json:"user_id,omitempty"`
    Published *bool   `json:"published,omitempty"`
    Search    *string `json:"search,omitempty"`
}

func (r *PostRepository) GetWithFilter(filter PostFilter, limit, offset int) ([]*models.Post, error) {
    query := r.db.Model(&models.Post{})

    if filter.UserID != nil {
        query = query.Where("user_id = ?", *filter.UserID)
    }
    if filter.Published != nil {
        query = query.Where("published = ?", *filter.Published)
    }
    if filter.Search != nil {
        query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+*filter.Search+"%", "%"+*filter.Search+"%")
    }

    var posts []*models.Post
    err := query.Limit(limit).Offset(offset).Find(&posts).Error
    return posts, err
}
```

This guide provides a complete template for adding new features while maintaining the clean architecture and consistency of your JWT authentication API.
