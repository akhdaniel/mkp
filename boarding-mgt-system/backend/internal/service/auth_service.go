package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/auth"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/ferryflow/boarding-mgt-system/internal/repository"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.LoginResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtUtil  *auth.JWTUtil
}

func NewAuthService(userRepo repository.UserRepository, jwtUtil *auth.JWTUtil) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Parse date of birth if provided
	var dob *time.Time
	if req.DateOfBirth != "" {
		parsedDob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format: %w", err)
		}
		dob = &parsedDob
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        &req.Phone,
		DateOfBirth:  dob,
		Nationality:  &req.Nationality,
		UserType:     "customer",
		IsVerified:   false,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Clear password hash from response
	user.PasswordHash = ""

	return user, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := auth.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate tokens
	accessToken, err := s.jwtUtil.GenerateToken(user.ID, user.Email, user.UserType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtUtil.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash tokens for storage
	accessTokenHash := auth.HashToken(accessToken)
	refreshTokenHash := auth.HashToken(refreshToken)

	// Create session
	session := &models.UserSession{
		UserID:           user.ID,
		TokenHash:        accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		IsActive:         true,
	}

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Non-critical error, log but don't fail
		fmt.Printf("failed to update last login: %v\n", err)
	}

	// Clear password hash from response
	user.PasswordHash = ""

	return &models.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    session.ExpiresAt,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.LoginResponse, error) {
	// Validate refresh token
	claims, err := s.jwtUtil.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate new tokens
	newAccessToken, err := s.jwtUtil.GenerateToken(user.ID, user.Email, user.UserType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtUtil.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash tokens for storage
	accessTokenHash := auth.HashToken(newAccessToken)
	refreshTokenHash := auth.HashToken(newRefreshToken)

	// Create new session
	session := &models.UserSession{
		UserID:           user.ID,
		TokenHash:        accessTokenHash,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		IsActive:         true,
	}

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Clear password hash from response
	user.PasswordHash = ""

	return &models.LoginResponse{
		User:         user,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    session.ExpiresAt,
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	// Deactivate all user sessions
	if err := s.userRepo.DeactivateUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	// Validate token
	claims, err := s.jwtUtil.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get session
	tokenHash := auth.HashToken(token)
	session, err := s.userRepo.GetSession(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired")
	}

	// Check if session is active
	if !session.IsActive {
		return nil, fmt.Errorf("session is not active")
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Clear password hash
	user.PasswordHash = ""

	return user, nil
}