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

type UserService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *models.ChangePasswordRequest) error
	ListUsers(ctx context.Context, filter *models.UserFilter) ([]*models.User, int, error)
	DeactivateUser(ctx context.Context, userID uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Clear password hash
	user.PasswordHash = ""

	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format: %w", err)
		}
		user.DateOfBirth = &dob
	}
	if req.Nationality != nil {
		user.Nationality = req.Nationality
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Clear password hash
	user.PasswordHash = ""

	return user, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, req *models.ChangePasswordRequest) error {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	if err := auth.VerifyPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Deactivate all sessions to force re-login
	if err := s.userRepo.DeactivateUserSessions(ctx, userID); err != nil {
		// Non-critical error, log but don't fail
		fmt.Printf("failed to deactivate sessions: %v\n", err)
	}

	return nil
}

func (s *userService) ListUsers(ctx context.Context, filter *models.UserFilter) ([]*models.User, int, error) {
	users, total, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	// Clear password hashes
	for _, user := range users {
		user.PasswordHash = ""
	}

	return users, total, nil
}

func (s *userService) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Deactivate user
	user.IsActive = false
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	// Deactivate all sessions
	if err := s.userRepo.DeactivateUserSessions(ctx, userID); err != nil {
		// Non-critical error, log but don't fail
		fmt.Printf("failed to deactivate sessions: %v\n", err)
	}

	return nil
}