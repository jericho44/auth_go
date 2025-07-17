package repositories

import (
	"errors"
	"time"

	"auth-jwt/models"

	"gorm.io/gorm"
)

// UserRepositoryInterface defines the contract for user data operations
type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	UpdateProfile(userID uint, email string) (*models.User, error)
	UpdatePassword(userID uint, hashedPassword string) error
	Delete(userID uint) error

	Exists(username, email string) (bool, error)
	GetUserStats(userID uint) (*models.UserStats, error)
	List(offset, limit int) ([]*models.User, int64, error)
	Search(query string, offset, limit int) ([]*models.User, int64, error)
}

// UserRepository implements UserRepositoryInterface using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

// Create inserts a new user into the database
func (ur *UserRepository) Create(user *models.User) error {
	result := ur.db.Create(user)
	return result.Error
}

// GetByID retrieves a user by their ID
func (ur *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := ur.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

// GetByUsername retrieves a user by their username
func (ur *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	result := ur.db.Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

// GetByEmail retrieves a user by their email
func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := ur.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

// Update updates a user's information
func (ur *UserRepository) Update(user *models.User) error {
	result := ur.db.Save(user)
	return result.Error
}

// UpdateProfile updates a user's profile information
func (ur *UserRepository) UpdateProfile(userID uint, email string) (*models.User, error) {
	var user models.User
	result := ur.db.Model(&user).Where("id = ?", userID).Update("email", email)
	if result.Error != nil {
		return nil, result.Error
	}

	// Fetch updated user
	if err := ur.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdatePassword updates a user's password
func (ur *UserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	result := ur.db.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword)
	return result.Error
}

// Delete soft deletes a user (GORM soft delete)
func (ur *UserRepository) Delete(userID uint) error {
	result := ur.db.Delete(&models.User{}, userID)
	return result.Error
}

// Exists checks if a user with the given username or email already exists
func (ur *UserRepository) Exists(username, email string) (bool, error) {
	var count int64
	result := ur.db.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	return count > 0, result.Error
}

// GetUserStats retrieves user statistics
func (ur *UserRepository) GetUserStats(userID uint) (*models.UserStats, error) {
	// Get basic user info
	user, err := ur.GetByID(userID)
	if err != nil || user == nil {
		return nil, err
	}

	stats := &models.UserStats{
		UserID:         user.ID,
		Username:       user.Username,
		Email:          user.Email,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		AccountAgeDays: int(time.Since(user.CreatedAt).Hours() / 24),
	}

	// Get login statistics
	var loginCount int64
	ur.db.Model(&models.LoginAttempt{}).Where("username = ? AND success = ?", user.Username, true).Count(&loginCount)
	stats.LoginCount = int(loginCount)

	// Get last successful login
	var lastLogin models.LoginAttempt
	if err := ur.db.Where("username = ? AND success = ?", user.Username, true).
		Order("attempted_at DESC").First(&lastLogin).Error; err == nil {
		stats.LastLoginAt = &lastLogin.AttemptedAt
	}

	return stats, nil
}

// List retrieves a paginated list of users
func (ur *UserRepository) List(offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// Get total count
	if err := ur.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	result := ur.db.Offset(offset).Limit(limit).Find(&users)
	return users, total, result.Error
}

// Search searches for users by username or email
func (ur *UserRepository) Search(query string, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	searchQuery := "%" + query + "%"
	condition := ur.db.Where("username ILIKE ? OR email ILIKE ?", searchQuery, searchQuery)

	// Get total count
	if err := condition.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	result := condition.Offset(offset).Limit(limit).Find(&users)
	return users, total, result.Error
}
