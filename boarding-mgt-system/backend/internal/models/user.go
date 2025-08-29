package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a system user
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	FirstName    string     `json:"first_name" db:"first_name"`
	LastName     string     `json:"last_name" db:"last_name"`
	Phone        *string    `json:"phone,omitempty" db:"phone"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
	Nationality  *string    `json:"nationality,omitempty" db:"nationality"`
	UserType     string     `json:"user_type" db:"user_type"`
	OperatorID   *uuid.UUID `json:"operator_id,omitempty" db:"operator_id"`
	IsVerified   bool       `json:"is_verified" db:"is_verified"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// UserSession represents an active user session
type UserSession struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash        string     `json:"-" db:"token_hash"`
	RefreshTokenHash string     `json:"-" db:"refresh_token_hash"`
	ExpiresAt        time.Time  `json:"expires_at" db:"expires_at"`
	IPAddress        *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent        *string    `json:"user_agent,omitempty" db:"user_agent"`
	IsActive         bool       `json:"is_active" db:"is_active"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// RegisterRequest represents user registration data
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Phone       string `json:"phone,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Nationality string `json:"nationality,omitempty"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	User         *User     `json:"user"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdateProfileRequest represents profile update data
type UpdateProfileRequest struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Nationality *string `json:"nationality,omitempty"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// UserFilter represents filters for listing users
type UserFilter struct {
	UserType   string     `json:"user_type,omitempty"`
	OperatorID *uuid.UUID `json:"operator_id,omitempty"`
	IsActive   *bool      `json:"is_active,omitempty"`
	Search     string     `json:"search,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Offset     int        `json:"offset,omitempty"`
}