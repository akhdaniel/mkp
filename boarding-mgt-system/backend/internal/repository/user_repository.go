package repository

import (
	"context"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *models.UserFilter) ([]*models.User, int, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	
	// Session management
	CreateSession(ctx context.Context, session *models.UserSession) error
	GetSession(ctx context.Context, tokenHash string) (*models.UserSession, error)
	DeactivateSession(ctx context.Context, sessionID uuid.UUID) error
	DeactivateUserSessions(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			email, password_hash, first_name, last_name, phone,
			date_of_birth, nationality, user_type, operator_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, is_verified, is_active, created_at, updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Phone,
		user.DateOfBirth, user.Nationality, user.UserType, user.OperatorID,
	).Scan(&user.ID, &user.IsVerified, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT 
			id, email, password_hash, first_name, last_name, phone,
			date_of_birth, nationality, user_type, operator_id,
			is_verified, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Phone,
		&user.DateOfBirth, &user.Nationality, &user.UserType, &user.OperatorID,
		&user.IsVerified, &user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT 
			id, email, password_hash, first_name, last_name, phone,
			date_of_birth, nationality, user_type, operator_id,
			is_verified, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	user := &models.User{}
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Phone,
		&user.DateOfBirth, &user.Nationality, &user.UserType, &user.OperatorID,
		&user.IsVerified, &user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET
			first_name = $2,
			last_name = $3,
			phone = $4,
			date_of_birth = $5,
			nationality = $6,
			is_verified = $7,
			is_active = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		user.ID, user.FirstName, user.LastName, user.Phone,
		user.DateOfBirth, user.Nationality, user.IsVerified, user.IsActive,
	).Scan(&user.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		return fmt.Errorf("user not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

func (r *userRepository) List(ctx context.Context, filter *models.UserFilter) ([]*models.User, int, error) {
	query := `
		SELECT 
			id, email, password_hash, first_name, last_name, phone,
			date_of_birth, nationality, user_type, operator_id,
			is_verified, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM users WHERE 1=1`
	
	args := []interface{}{}
	argCount := 0
	
	// Build filters
	if filter.UserType != "" {
		argCount++
		query += fmt.Sprintf(" AND user_type = $%d", argCount)
		countQuery += fmt.Sprintf(" AND user_type = $%d", argCount)
		args = append(args, filter.UserType)
	}
	
	if filter.OperatorID != nil {
		argCount++
		query += fmt.Sprintf(" AND operator_id = $%d", argCount)
		countQuery += fmt.Sprintf(" AND operator_id = $%d", argCount)
		args = append(args, *filter.OperatorID)
	}
	
	if filter.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		countQuery += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filter.IsActive)
	}
	
	if filter.Search != "" {
		argCount++
		searchPattern := "%" + filter.Search + "%"
		query += fmt.Sprintf(" AND (email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)", argCount, argCount, argCount)
		countQuery += fmt.Sprintf(" AND (email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, searchPattern)
	}
	
	// Get total count
	var totalCount int
	err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}
	
	// Add pagination
	query += " ORDER BY created_at DESC"
	if filter.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}
	
	// Execute query
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()
	
	users := []*models.User{}
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Phone,
			&user.DateOfBirth, &user.Nationality, &user.UserType, &user.OperatorID,
			&user.IsVerified, &user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	return users, totalCount, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = $1`
	
	_, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	
	return nil
}

func (r *userRepository) CreateSession(ctx context.Context, session *models.UserSession) error {
	query := `
		INSERT INTO user_sessions (
			user_id, token_hash, refresh_token_hash, expires_at, ip_address, user_agent
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		session.UserID, session.TokenHash, session.RefreshTokenHash,
		session.ExpiresAt, session.IPAddress, session.UserAgent,
	).Scan(&session.ID, &session.CreatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	
	session.IsActive = true
	return nil
}

func (r *userRepository) GetSession(ctx context.Context, tokenHash string) (*models.UserSession, error) {
	query := `
		SELECT 
			id, user_id, token_hash, refresh_token_hash, expires_at,
			ip_address, user_agent, is_active, created_at
		FROM user_sessions
		WHERE token_hash = $1 AND is_active = true AND expires_at > CURRENT_TIMESTAMP
	`
	
	session := &models.UserSession{}
	err := r.db.Pool.QueryRow(ctx, query, tokenHash).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.RefreshTokenHash,
		&session.ExpiresAt, &session.IPAddress, &session.UserAgent,
		&session.IsActive, &session.CreatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("session not found or expired")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	
	return session, nil
}

func (r *userRepository) DeactivateSession(ctx context.Context, sessionID uuid.UUID) error {
	query := `UPDATE user_sessions SET is_active = false WHERE id = $1`
	
	_, err := r.db.Pool.Exec(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to deactivate session: %w", err)
	}
	
	return nil
}

func (r *userRepository) DeactivateUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE user_sessions SET is_active = false WHERE user_id = $1 AND is_active = true`
	
	_, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user sessions: %w", err)
	}
	
	return nil
}