package services

import (
	"errors"

	"auth-jwt/models"
	"auth-jwt/repositories"
)

// UserServiceInterface defines the contract for user business operations
type UserServiceInterface interface {
	CreateUser(username, email, password string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUserProfile(userID int, email string) (*models.User, error)
	UpdateUserPassword(userID int, newPassword string) error
	DeleteUser(userID int) error
	GetUserStats(userID int) (*models.UserStats, error)
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

	// Create user model
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	// Save to database
	if err := us.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (us *UserService) GetUserByID(id int) (*models.User, error) {
	return us.userRepo.GetByID(id)
}

// GetUserByUsername retrieves a user by username
func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	return us.userRepo.GetByUsername(username)
}

// UpdateUserProfile updates user profile information
func (us *UserService) UpdateUserProfile(userID int, email string) (*models.User, error) {
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
func (us *UserService) UpdateUserPassword(userID int, newPassword string) error {
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

// DeleteUser deletes a user account
func (us *UserService) DeleteUser(userID int) error {
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

// GetUserStats retrieves user statistics
func (us *UserService) GetUserStats(userID int) (*models.UserStats, error) {
	return us.userRepo.GetUserStats(userID)
}
