package services

import (
	"errors"

	"auth-jwt/models"
	"auth-jwt/repositories"
)

// UserServiceInterface defines the contract for user business operations
type UserServiceInterface interface {
	CreateUser(username, email, password string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUserProfile(userID uint, email string) (*models.User, error)
	UpdateUserPassword(userID uint, newPassword string) error
	DeleteUser(userID uint) error
	SoftDeleteUser(userID uint) error
	GetUserStats(userID uint) (*models.UserStats, error)
	ListUsers(page, limit int) ([]*models.User, int64, error)
	SearchUsers(query string, page, limit int) ([]*models.User, int64, error)
}

// UserService implements UserServiceInterface
type UserService struct {
	userRepo repositories.UserRepositoryInterface
}

// NewUserService creates a new user service instance
func NewUserService(userRepo repositories.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user with validation
func (us *UserService) CreateUser(username, email, password string) (*models.User, error) {
	// Check if user already exists
	exists, err := us.userRepo.Exists(username, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Create user model (password will be hashed by GORM hook)
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Save to database
	if err := us.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (us *UserService) GetUserByID(id uint) (*models.User, error) {
	return us.userRepo.GetByID(id)
}

// GetUserByUsername retrieves a user by username
func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	return us.userRepo.GetByUsername(username)
}

// UpdateUserProfile updates user profile information
func (us *UserService) UpdateUserProfile(userID uint, email string) (*models.User, error) {
	// Check if user exists
	user, err := us.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update profile
	return us.userRepo.UpdateProfile(userID, email)
}

// UpdateUserPassword updates user password
func (us *UserService) UpdateUserPassword(userID uint, newPassword string) error {
	// Check if user exists
	user, err := us.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Hash new password
	tempUser := &models.User{Password: newPassword}
	if err := tempUser.HashPassword(); err != nil {
		return err
	}

	// Update password
	return us.userRepo.UpdatePassword(userID, tempUser.Password)
}

// DeleteUser permanently deletes a user account
func (us *UserService) DeleteUser(userID uint) error {
	// Check if user exists
	user, err := us.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Delete user
	return us.userRepo.Delete(userID)
}

// SoftDeleteUser soft deletes a user account (GORM handles soft delete automatically)
func (us *UserService) SoftDeleteUser(userID uint) error {
	// Check if user exists
	user, err := us.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// GORM Delete method performs soft delete automatically when DeletedAt field is present
	return us.userRepo.Delete(userID)
}

// GetUserStats retrieves user statistics
func (us *UserService) GetUserStats(userID uint) (*models.UserStats, error) {
	return us.userRepo.GetUserStats(userID)
}

// ListUsers retrieves a paginated list of users
func (us *UserService) ListUsers(page, limit int) ([]*models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return us.userRepo.List(offset, limit)
}

// SearchUsers searches for users by username or email
func (us *UserService) SearchUsers(query string, page, limit int) ([]*models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	return us.userRepo.Search(query, offset, limit)
}
