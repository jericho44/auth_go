# Adding New Features Guide

This guide explains how to add new features to the JWT Authentication API following the established clean architecture patterns and transaction support.

## Architecture Overview

The project follows clean architecture with these layers:

1. **Models** - Data structures and domain entities
2. **Repositories** - Data access layer with transaction support
3. **Services** - Business logic layer
4. **Controllers** - HTTP request/response handling
5. **Handlers** - Route handlers
6. **Routes** - Route definitions

## Step-by-Step Feature Addition

### 1. Define the Model

Create or update models in `models/` directory:

```go
// models/example.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Example struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"not null"`
    UserID    uint           `json:"user_id" gorm:"not null"`
    User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
    IsActive  bool           `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete support
}

// Request/Response DTOs
type CreateExampleRequest struct {
    Name string `json:"name" validate:"required,min=3,max=100"`
}

type UpdateExampleRequest struct {
    Name     *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
    IsActive *bool   `json:"is_active,omitempty"`
}

type ExampleResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    UserID    uint      `json:"user_id"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. Create Repository Interface and Implementation

Create repository in `repositories/` directory:

```go
// repositories/example_repository.go
package repositories

import (
    "auth-jwt/models"
    "gorm.io/gorm"
)

type ExampleRepositoryInterface interface {
    Create(example *models.Example) error
    GetByID(id uint) (*models.Example, error)
    GetByUserID(userID uint, limit, offset int) ([]*models.Example, error)
    Update(example *models.Example) error
    Delete(id uint) error
    CountByUserID(userID uint) (int64, error)
}

type ExampleRepository struct {
    db *gorm.DB
}

func NewExampleRepository(db *gorm.DB) ExampleRepositoryInterface {
    return &ExampleRepository{db: db}
}

// Create with transaction support
func (r *ExampleRepository) Create(example *models.Example) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        return tx.Create(example).Error
    })
}

// GetByID - read operations don't need transactions
func (r *ExampleRepository) GetByID(id uint) (*models.Example, error) {
    var example models.Example
    err := r.db.Preload("User").First(&example, id).Error
    if err != nil {
        return nil, err
    }
    return &example, nil
}

func (r *ExampleRepository) GetByUserID(userID uint, limit, offset int) ([]*models.Example, error) {
    var examples []*models.Example
    err := r.db.Where("user_id = ?", userID).
        Limit(limit).Offset(offset).
        Order("created_at DESC").
        Find(&examples).Error
    return examples, err
}

// Update with transaction support
func (r *ExampleRepository) Update(example *models.Example) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        return tx.Save(example).Error
    })
}

// Delete with transaction support
func (r *ExampleRepository) Delete(id uint) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        return tx.Delete(&models.Example{}, id).Error
    })
}

func (r *ExampleRepository) CountByUserID(userID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.Example{}).Where("user_id = ?", userID).Count(&count).Error
    return count, err
}
```

### 3. Create Service Layer

Create service in `services/` directory:

```go
// services/example_service.go
package services

import (
    "errors"
    "auth-jwt/models"
    "auth-jwt/repositories"
)

type ExampleServiceInterface interface {
    CreateExample(userID uint, req models.CreateExampleRequest) (*models.ExampleResponse, error)
    GetExample(id, userID uint) (*models.ExampleResponse, error)
    GetUserExamples(userID uint, page, limit int) ([]*models.ExampleResponse, int64, error)
    UpdateExample(id, userID uint, req models.UpdateExampleRequest) (*models.ExampleResponse, error)
    DeleteExample(id, userID uint) error
}

type ExampleService struct {
    exampleRepo repositories.ExampleRepositoryInterface
    userRepo    repositories.UserRepositoryInterface
}

func NewExampleService(exampleRepo repositories.ExampleRepositoryInterface, userRepo repositories.UserRepositoryInterface) ExampleServiceInterface {
    return &ExampleService{
        exampleRepo: exampleRepo,
        userRepo:    userRepo,
    }
}

func (s *ExampleService) CreateExample(userID uint, req models.CreateExampleRequest) (*models.ExampleResponse, error) {
    // Verify user exists
    user, err := s.userRepo.GetByID(userID)
    if err != nil || user == nil {
        return nil, errors.New("user not found")
    }

    // Create example
    example := &models.Example{
        Name:   req.Name,
        UserID: userID,
    }

    // Repository handles transaction
    if err := s.exampleRepo.Create(example); err != nil {
        return nil, errors.New("failed to create example")
    }

    return s.convertToResponse(example), nil
}

func (s *ExampleService) GetExample(id, userID uint) (*models.ExampleResponse, error) {
    example, err := s.exampleRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("example not found")
    }

    // Check ownership
    if example.UserID != userID {
        return nil, errors.New("unauthorized")
    }

    return s.convertToResponse(example), nil
}

func (s *ExampleService) GetUserExamples(userID uint, page, limit int) ([]*models.ExampleResponse, int64, error) {
    offset := (page - 1) * limit
    examples, err := s.exampleRepo.GetByUserID(userID, limit, offset)
    if err != nil {
        return nil, 0, errors.New("failed to retrieve examples")
    }

    total, err := s.exampleRepo.CountByUserID(userID)
    if err != nil {
        return nil, 0, errors.New("failed to count examples")
    }

    var responses []*models.ExampleResponse
    for _, example := range examples {
        responses = append(responses, s.convertToResponse(example))
    }

    return responses, total, nil
}

func (s *ExampleService) UpdateExample(id, userID uint, req models.UpdateExampleRequest) (*models.ExampleResponse, error) {
    example, err := s.exampleRepo.GetByID(id)
    if err != nil {
        return nil, errors.New("example not found")
    }

    // Check ownership
    if example.UserID != userID {
        return nil, errors.New("unauthorized")
    }

    // Update fields
    if req.Name != nil {
        example.Name = *req.Name
    }
    if req.IsActive != nil {
        example.IsActive = *req.IsActive
    }

    // Repository handles transaction
    if err := s.exampleRepo.Update(example); err != nil {
        return nil, errors.New("failed to update example")
    }

    return s.convertToResponse(example), nil
}

func (s *ExampleService) DeleteExample(id, userID uint) error {
    example, err := s.exampleRepo.GetByID(id)
    if err != nil {
        return errors.New("example not found")
    }

    // Check ownership
    if example.UserID != userID {
        return errors.New("unauthorized")
    }

    // Repository handles transaction
    return s.exampleRepo.Delete(id)
}

// Helper method
func (s *ExampleService) convertToResponse(example *models.Example) *models.ExampleResponse {
    return &models.ExampleResponse{
        ID:        example.ID,
        Name:      example.Name,
        UserID:    example.UserID,
        IsActive:  example.IsActive,
        CreatedAt: example.CreatedAt,
        UpdatedAt: example.UpdatedAt,
    }
}
```

### 4. Create Controller

Create controller in `controllers/` directory:

```go
// controllers/example_controller.go
package controllers

import (
    "net/http"
    "strconv"

    "auth-jwt/models"
    "auth-jwt/services"
    "auth-jwt/utils"

    "github.com/gin-gonic/gin"
)

type ExampleController struct {
    exampleService services.ExampleServiceInterface
}

func NewExampleController(exampleService services.ExampleServiceInterface) *ExampleController {
    return &ExampleController{
        exampleService: exampleService,
    }
}

// @Summary Create example
// @Description Create a new example
// @Tags examples
// @Accept json
// @Produce json
// @Param request body models.CreateExampleRequest true "Example data"
// @Success 201 {object} utils.Response{data=models.ExampleResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Security BearerAuth
// @Router /api/examples [post]
func (ctrl *ExampleController) CreateExample(c *gin.Context) {
    userID := c.GetUint("user_id")

    var req models.CreateExampleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
        return
    }

    // Validate request
    if err := utils.ValidateStruct(req); err != nil {
        utils.ValidationErrorResponse(c, err)
        return
    }

    example, err := ctrl.exampleService.CreateExample(userID, req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
        return
    }

    utils.SuccessResponse(c, http.StatusCreated, "Example created successfully", example)
}

// @Summary Get example
// @Description Get example by ID
// @Tags examples
// @Produce json
// @Param id path int true "Example ID"
// @Success 200 {object} utils.Response{data=models.ExampleResponse}
// @Failure 404 {object} utils.Response
// @Security BearerAuth
// @Router /api/examples/{id} [get]
func (ctrl *ExampleController) GetExample(c *gin.Context) {
    userID := c.GetUint("user_id")

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid example ID", nil)
        return
    }

    example, err := ctrl.exampleService.GetExample(uint(id), userID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Example retrieved successfully", example)
}

// @Summary List user examples
// @Description Get paginated list of user's examples
// @Tags examples
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response{data=object}
// @Security BearerAuth
// @Router /api/examples [get]
func (ctrl *ExampleController) GetUserExamples(c *gin.Context) {
    userID := c.GetUint("user_id")

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 10
    }

    examples, total, err := ctrl.exampleService.GetUserExamples(userID, page, limit)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
        return
    }

    response := map[string]interface{}{
        "examples": examples,
        "total":    total,
        "page":     page,
        "limit":    limit,
    }

    utils.SuccessResponse(c, http.StatusOK, "Examples retrieved successfully", response)
}

// @Summary Update example
// @Description Update example by ID
// @Tags examples
// @Accept json
// @Produce json
// @Param id path int true "Example ID"
// @Param request body models.UpdateExampleRequest true "Update data"
// @Success 200 {object} utils.Response{data=models.ExampleResponse}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Security BearerAuth
// @Router /api/examples/{id} [put]
func (ctrl *ExampleController) UpdateExample(c *gin.Context) {
    userID := c.GetUint("user_id")

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid example ID", nil)
        return
    }

    var req models.UpdateExampleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
        return
    }

    // Validate request
    if err := utils.ValidateStruct(req); err != nil {
        utils.ValidationErrorResponse(c, err)
        return
    }

    example, err := ctrl.exampleService.UpdateExample(uint(id), userID, req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Example updated successfully", example)
}

// @Summary Delete example
// @Description Delete example by ID
// @Tags examples
// @Param id path int true "Example ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Security BearerAuth
// @Router /api/examples/{id} [delete]
func (ctrl *ExampleController) DeleteExample(c *gin.Context) {
    userID := c.GetUint("user_id")

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid example ID", nil)
        return
    }

    if err := ctrl.exampleService.DeleteExample(uint(id), userID); err != nil {
        utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
        return
    }

    utils.SuccessResponse(c, http.StatusOK, "Example deleted successfully", nil)
}
```

### 5. Create Routes

Create routes in `routes/` directory:

```go
// routes/example.go
package routes

import (
    "auth-jwt/controllers"
    "auth-jwt/middleware"

    "github.com/gin-gonic/gin"
)

func SetupExampleRoutes(router *gin.RouterGroup, exampleController *controllers.ExampleController) {
    examples := router.Group("/examples")
    examples.Use(middleware.AuthMiddleware()) // Require authentication
    {
        examples.POST("", exampleController.CreateExample)
        examples.GET("", exampleController.GetUserExamples)
        examples.GET("/:id", exampleController.GetExample)
        examples.PUT("/:id", exampleController.UpdateExample)
        examples.DELETE("/:id", exampleController.DeleteExample)
    }
}
```

### 6. Update Main Application

Update `main.go` to wire up the new feature:

```go
// Add to main.go
func main() {
    // ... existing code ...

    // Initialize repositories
    exampleRepo := repositories.NewExampleRepository(database.GetDB())

    // Initialize services
    exampleService := services.NewExampleService(exampleRepo, userRepo)

    // Initialize controllers
    exampleController := controllers.NewExampleController(exampleService)

    // Setup routes
    api := router.Group("/api")
    routes.SetupExampleRoutes(api, exampleController)

    // ... rest of the code ...
}
```

### 7. Update Database Migration

Add the new model to `database/db.go`:

```go
// Update AutoMigrate function
func AutoMigrate() error {
    err := DB.AutoMigrate(
        &models.User{},
        &models.PasswordResetToken{},
        &models.LoginAttempt{},
        &models.File{},
        &models.Example{}, // Add new model
    )

    if err != nil {
        return fmt.Errorf("failed to auto-migrate: %w", err)
    }

    return nil
}
```

## Transaction Best Practices

### When to Use Transactions

1. **Always use transactions for write operations** (Create, Update, Delete)
2. **Use service-level transactions** for operations spanning multiple repositories
3. **Keep read operations outside transactions** for better performance
4. **Validate data before starting transactions** to avoid unnecessary rollbacks
5. **Use transactions for operations that must be atomic** (all succeed or all fail)

### Transaction Patterns

#### Repository Level (Simple Operations)

```go
func (r *Repository) Create(entity *Model) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        return tx.Create(entity).Error
    })
}

func (r *Repository) Update(entity *Model) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        return tx.Save(entity).Error
    })
}

func (r *Repository) Delete(id uint) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        return tx.Delete(&Model{}, id).Error
    })
}
```

#### Service Level (Complex Operations)

```go
func (s *Service) ComplexOperation(data Data) error {
    // Validate outside transaction
    if err := s.validateData(data); err != nil {
        return err
    }

    db := s.getDB()
    return db.Transaction(func(tx *gorm.DB) error {
        // Multiple operations that must succeed together
        if err := tx.Create(&entity1).Error; err != nil {
            return err
        }

        if err := tx.Model(&entity2).Where("id = ?", data.ID).Update("status", "updated").Error; err != nil {
            return err
        }

        // Delete related records
        if err := tx.Where("parent_id = ?", entity1.ID).Delete(&RelatedModel{}).Error; err != nil {
            return err
        }

        return nil
    })
}

// Helper method to get DB instance
func (s *Service) getDB() *gorm.DB {
    return database.GetDB()
}
```

#### Advanced Transaction Patterns

**File Operations with Database Consistency:**

```go
func (s *FileService) UploadWithMetadata(file FileData, metadata Metadata) error {
    // Save file first
    filePath, err := s.saveFile(file)
    if err != nil {
        return err
    }

    // Use transaction for database operations
    db := s.getDB()
    err = db.Transaction(func(tx *gorm.DB) error {
        // Create file record
        fileRecord := &models.File{
            Path: filePath,
            Name: file.Name,
        }
        if err := tx.Create(fileRecord).Error; err != nil {
            return err
        }

        // Create metadata record
        metadataRecord := &models.FileMetadata{
            FileID: fileRecord.ID,
            Data:   metadata,
        }
        if err := tx.Create(metadataRecord).Error; err != nil {
            return err
        }

        return nil
    })

    // Clean up file if database transaction failed
    if err != nil {
        os.Remove(filePath)
        return err
    }

    return nil
}
```

**Batch Operations with Transaction:**

```go
func (s *Service) BatchUpdate(items []UpdateItem) error {
    if len(items) == 0 {
        return nil
    }

    db := s.getDB()
    return db.Transaction(func(tx *gorm.DB) error {
        for _, item := range items {
            if err := tx.Model(&Model{}).Where("id = ?", item.ID).Updates(item.Data).Error; err != nil {
                return err // This will rollback all previous updates
            }
        }
        return nil
    })
}
```

### Transaction Error Handling

```go
func (s *Service) SafeOperation(data Data) error {
    db := s.getDB()
    return db.Transaction(func(tx *gorm.DB) error {
        // Operation 1
        if err := tx.Create(&entity1).Error; err != nil {
            // Log specific error
            log.Printf("Failed to create entity1: %v", err)
            return fmt.Errorf("failed to create primary record: %w", err)
        }

        // Operation 2
        if err := tx.Create(&entity2).Error; err != nil {
            log.Printf("Failed to create entity2: %v", err)
            return fmt.Errorf("failed to create secondary record: %w", err)
        }

        // Validation after operations
        if !s.validateBusinessRules(entity1, entity2) {
            return errors.New("business rule validation failed")
        }

        return nil
    })
}
```

### Performance Considerations

1. **Keep transactions short** - Long transactions hold locks longer
2. **Validate before transactions** - Avoid rollbacks when possible
3. **Use read replicas for queries** - Don't include reads in write transactions
4. **Batch operations efficiently** - Group related operations

```go
// Good: Validate first, then transact
func (s *Service) EfficientCreate(data Data) error {
    // Validate outside transaction
    if err := s.validate(data); err != nil {
        return err
    }

    // Check business rules outside transaction
    if exists, err := s.checkExists(data.Key); err != nil {
        return err
    } else if exists {
        return errors.New("already exists")
    }

    // Short transaction for actual write
    return s.repo.Create(&Model{Data: data})
}

// Bad: Long transaction with validation inside
func (s *Service) InefficientCreate(data Data) error {
    db := s.getDB()
    return db.Transaction(func(tx *gorm.DB) error {
        // This holds locks while validating
        if err := s.validate(data); err != nil {
            return err
        }

        var existing Model
        if err := tx.Where("key = ?", data.Key).First(&existing).Error; err == nil {
            return errors.New("already exists")
        }

        return tx.Create(&Model{Data: data}).Error
    })
}
```

## Testing Your Feature

### Unit Tests

Create tests in `*_test.go` files:

```go
// services/example_service_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestExampleService_CreateExample(t *testing.T) {
    // Test implementation
}
```

### Integration Tests

Test the complete flow from HTTP request to database:

```go
func TestExampleAPI_CreateExample(t *testing.T) {
    // Setup test database
    // Make HTTP request
    // Assert response
    // Verify database state
}
```

## Documentation

1. **Update Swagger comments** for API documentation
2. **Add examples to API documentation**
3. **Update README.md** with new endpoints
4. **Create feature-specific documentation** if needed

## Checklist

- [ ] Model created with proper validation tags
- [ ] Repository interface and implementation with transactions
- [ ] Service layer with business logic
- [ ] Controller with proper error handling
- [ ] Routes with authentication middleware
- [ ] Database migration updated
- [ ] Main application wired up
- [ ] Swagger documentation added
- [ ] Tests written
- [ ] README updated
- [ ] API examples documented

Following this guide ensures your new features integrate seamlessly with the existing architecture and maintain the same quality standards.
