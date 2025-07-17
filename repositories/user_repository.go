package repositories

import (
	"database/sql"
	"time"

	"auth-jwt/models"
)

// UserRepositoryInterface defines the contract for user data operations
type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	UpdateProfile(userID int, email string) (*models.User, error)
	UpdatePassword(userID int, hashedPassword string) error
	Delete(userID int) error
	Exists(username, email string) (bool, error)
	GetUserStats(userID int) (*models.UserStats, error)
}

// UserRepository implements UserRepositoryInterface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

// Create inserts a new user into the database
func (ur *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := ur.db.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

// GetByID retrieves a user by their ID
func (ur *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1`
	err := ur.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetByUsername retrieves a user by their username
func (ur *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1`
	err := ur.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetByEmail retrieves a user by their email
func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`
	err := ur.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Update updates a user's information
func (ur *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET username = $1, email = $2, password = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := ur.db.Exec(query, user.Username, user.Email, user.Password, user.ID)
	return err
}

// UpdateProfile updates a user's profile information
func (ur *UserRepository) UpdateProfile(userID int, email string) (*models.User, error) {
	query := `UPDATE users SET email = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING id, username, email, created_at, updated_at`
	user := &models.User{}
	err := ur.db.QueryRow(query, email, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdatePassword updates a user's password
func (ur *UserRepository) UpdatePassword(userID int, hashedPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := ur.db.Exec(query, hashedPassword, userID)
	return err
}

// Delete removes a user from the database
func (ur *UserRepository) Delete(userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := ur.db.Exec(query, userID)
	return err
}

// Exists checks if a user with the given username or email already exists
func (ur *UserRepository) Exists(username, email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2`
	err := ur.db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserStats retrieves user statistics
func (ur *UserRepository) GetUserStats(userID int) (*models.UserStats, error) {
	stats := &models.UserStats{}

	// Get basic user info
	user, err := ur.GetByID(userID)
	if err != nil || user == nil {
		return nil, err
	}

	stats.UserID = user.ID
	stats.Username = user.Username
	stats.Email = user.Email
	stats.CreatedAt = user.CreatedAt
	stats.UpdatedAt = user.UpdatedAt

	// Calculate account age in days
	stats.AccountAgeDays = int(time.Since(user.CreatedAt).Hours() / 24)

	// TODO: Add more statistics like login count, last login, etc.
	// This would require additional tables for tracking user activity

	return stats, nil
}
