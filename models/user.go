package models

import (
	"database/sql"
	"time"

	"auth-jwt/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't include in JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// CreateUser creates a new user in database
func CreateUser(username, email, password string) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := database.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves user by username
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE username = $1`
	err := database.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves user by ID
func GetUserByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// UpdateUserProfile updates user profile information
func UpdateUserProfile(userID int, email string) (*User, error) {
	query := `UPDATE users SET email = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING id, username, email, created_at, updated_at`
	user := &User{}
	err := database.DB.QueryRow(query, email, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUserPassword updates user password
func UpdateUserPassword(userID int, newPassword string) error {
	// Hash the new password
	user := &User{Password: newPassword}
	if err := user.HashPassword(); err != nil {
		return err
	}

	query := `UPDATE users SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := database.DB.Exec(query, user.Password, userID)
	return err
}

// DeleteUser deletes a user account
func DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := database.DB.Exec(query, userID)
	return err
}
