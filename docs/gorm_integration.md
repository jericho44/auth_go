# GORM Integration Guide

This document explains how GORM (Go Object-Relational Mapping) has been integrated into the authentication project.

## Overview

GORM is a powerful ORM library for Go that provides:

- Auto-migration capabilities
- Model associations and relationships
- Query builder with method chaining
- Hooks for lifecycle events
- Soft delete support
- Connection pooling and performance optimization

## Models

### User Model

```go
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
    Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
    Password  string         `json:"-" gorm:"size:255;not null"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete
}
```

### Password Reset Token Model

```go
type PasswordResetToken struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    UserID    uint           `json:"user_id" gorm:"not null;index"`
    User      User           `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
    Token     string         `json:"token" gorm:"uniqueIndex;size:255;not null"`
    ExpiresAt time.Time      `json:"expires_at" gorm:"not null;index"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### Login Attempt Model

```go
type LoginAttempt struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Username    string         `json:"username" gorm:"size:50;not null;index"`
    IPAddress   string         `json:"ip_address" gorm:"size:45"`
    Success     bool           `json:"success" gorm:"not null;default:false"`
    AttemptedAt time.Time      `json:"attempted_at" gorm:"not null;index"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
```

## GORM Features Used

### 1. Auto-Migration

```go
func AutoMigrate() error {
    return DB.AutoMigrate(
        &models.User{},
        &models.PasswordResetToken{},
        &models.LoginAttempt{},
    )
}
```

### 2. Model Hooks

```go
// BeforeCreate hook for password hashing
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.Password != "" {
        return u.HashPassword()
    }
    return nil
}

// BeforeUpdate hook for password hashing
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    if tx.Statement.Changed("Password") && u.Password != "" {
        return u.HashPassword()
    }
    return nil
}
```

### 3. Soft Delete

- Users are soft deleted by default (DeletedAt field)
- Allows data recovery and audit trails
- Automatically excluded from queries unless explicitly included

### 4. Associations

- Foreign key relationships between models
- Cascade delete constraints
- Preloading support for related data

## Repository Pattern with GORM

### User Repository Example

```go
// Create a new user
func (ur *UserRepository) Create(user *models.User) error {
    return ur.db.Create(user).Error
}

// Get user by ID
func (ur *UserRepository) GetByID(id uint) (*models.User, error) {
    var user models.User
    err := ur.db.First(&user, id).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

// Update user profile
func (ur *UserRepository) UpdateProfile(userID uint, email string) (*models.User, error) {
    var user models.User
    err := ur.db.Model(&user).Where("id = ?", userID).Update("email", email).Error
    if err != nil {
        return nil, err
    }

    // Fetch updated user
    err = ur.db.First(&user, userID).Error
    if err != nil {
        return nil, err
    }

    return &user, nil
}
```

## Query Examples

### Basic Queries

```go
// Find by ID
var user User
db.First(&user, 1)

// Find by condition
var user User
db.Where("username = ?", "john").First(&user)

// Find multiple records
var users []User
db.Where("created_at > ?", time.Now().AddDate(0, -1, 0)).Find(&users)
```

### Advanced Queries

```go
// Count records
var count int64
db.Model(&User{}).Where("created_at > ?", time.Now().AddDate(0, -1, 0)).Count(&count)

// Order and limit
var users []User
db.Order("created_at DESC").Limit(10).Find(&users)

// Select specific fields
var users []User
db.Select("id, username, email").Find(&users)
```

### Aggregations

```go
// Get user statistics
var loginCount int64
db.Model(&LoginAttempt{}).Where("username = ? AND success = ?", username, true).Count(&loginCount)

// Get latest login
var lastLogin LoginAttempt
db.Where("username = ? AND success = ?", username, true).
   Order("attempted_at DESC").
   First(&lastLogin)
```

## Performance Optimizations

### 1. Indexes

- Unique indexes on username and email
- Regular indexes on frequently queried fields
- Composite indexes for complex queries

### 2. Connection Pooling

```go
sqlDB, err := db.DB()
if err != nil {
    return err
}

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
sqlDB.SetMaxIdleConns(10)

// SetMaxOpenConns sets the maximum number of open connections to the database
sqlDB.SetMaxOpenConns(100)

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 3. Preloading

```go
// Preload related data
var tokens []PasswordResetToken
db.Preload("User").Find(&tokens)
```

## Migration from Raw SQL

### Before (Raw SQL)

```go
query := `SELECT id, username, email FROM users WHERE username = $1`
err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email)
```

### After (GORM)

```go
err := db.Where("username = ?", username).First(&user).Error
```

## Benefits of GORM Integration

### 1. **Reduced Boilerplate**

- No manual SQL query writing
- Automatic struct scanning
- Built-in error handling

### 2. **Type Safety**

- Compile-time checking
- No SQL injection vulnerabilities
- Automatic data type conversion

### 3. **Maintainability**

- Auto-migration keeps schema in sync
- Model-driven development
- Easy to refactor and extend

### 4. **Performance**

- Connection pooling
- Query optimization
- Lazy loading support

### 5. **Features**

- Soft delete support
- Model hooks and callbacks
- Transaction support
- Association handling

## Best Practices

### 1. **Model Design**

- Use appropriate GORM tags
- Define relationships clearly
- Implement necessary hooks

### 2. **Query Optimization**

- Use indexes appropriately
- Avoid N+1 queries with preloading
- Use raw SQL for complex queries when needed

### 3. **Error Handling**

- Check for `gorm.ErrRecordNotFound`
- Handle database connection errors
- Use transactions for data consistency

### 4. **Security**

- GORM prevents SQL injection by default
- Validate input data before database operations
- Use proper access controls

## Testing with GORM

### 1. **Test Database Setup**

```go
func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Auto-migrate test schema
    db.AutoMigrate(&User{}, &PasswordResetToken{}, &LoginAttempt{})

    return db
}
```

### 2. **Repository Testing**

```go
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB()
    repo := NewUserRepository(db)

    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }

    err := repo.Create(user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

This GORM integration provides a robust, maintainable, and performant data access layer for the authentication system.
