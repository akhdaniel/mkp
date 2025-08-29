package repository

import (
	"context"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OperatorRepository interface {
	Create(ctx context.Context, operator *models.Operator) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Operator, error)
	GetByCode(ctx context.Context, code string) (*models.Operator, error)
	Update(ctx context.Context, operator *models.Operator) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Operator, int, error)
}

type operatorRepository struct {
	db *database.DB
}

func NewOperatorRepository(db *database.DB) OperatorRepository {
	return &operatorRepository{db: db}
}

func (r *operatorRepository) Create(ctx context.Context, operator *models.Operator) error {
	query := `
		INSERT INTO operators (
			name, code, contact_email, contact_phone, address, settings
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, is_active, created_at, updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		operator.Name, operator.Code, operator.ContactEmail,
		operator.ContactPhone, operator.Address, operator.Settings,
	).Scan(&operator.ID, &operator.IsActive, &operator.CreatedAt, &operator.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create operator: %w", err)
	}
	
	return nil
}

func (r *operatorRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Operator, error) {
	query := `
		SELECT 
			id, name, code, contact_email, contact_phone, address,
			is_active, settings, created_at, updated_at
		FROM operators
		WHERE id = $1
	`
	
	operator := &models.Operator{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&operator.ID, &operator.Name, &operator.Code, &operator.ContactEmail,
		&operator.ContactPhone, &operator.Address, &operator.IsActive,
		&operator.Settings, &operator.CreatedAt, &operator.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("operator not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get operator: %w", err)
	}
	
	return operator, nil
}

func (r *operatorRepository) GetByCode(ctx context.Context, code string) (*models.Operator, error) {
	query := `
		SELECT 
			id, name, code, contact_email, contact_phone, address,
			is_active, settings, created_at, updated_at
		FROM operators
		WHERE code = $1
	`
	
	operator := &models.Operator{}
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(
		&operator.ID, &operator.Name, &operator.Code, &operator.ContactEmail,
		&operator.ContactPhone, &operator.Address, &operator.IsActive,
		&operator.Settings, &operator.CreatedAt, &operator.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("operator not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get operator: %w", err)
	}
	
	return operator, nil
}

func (r *operatorRepository) Update(ctx context.Context, operator *models.Operator) error {
	query := `
		UPDATE operators SET
			name = $2,
			contact_email = $3,
			contact_phone = $4,
			address = $5,
			is_active = $6,
			settings = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		operator.ID, operator.Name, operator.ContactEmail,
		operator.ContactPhone, operator.Address, operator.IsActive, operator.Settings,
	).Scan(&operator.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		return fmt.Errorf("operator not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update operator: %w", err)
	}
	
	return nil
}

func (r *operatorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM operators WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete operator: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("operator not found")
	}
	
	return nil
}

func (r *operatorRepository) List(ctx context.Context, limit, offset int) ([]*models.Operator, int, error) {
	query := `
		SELECT 
			id, name, code, contact_email, contact_phone, address,
			is_active, settings, created_at, updated_at
		FROM operators
		ORDER BY created_at DESC
	`
	countQuery := `SELECT COUNT(*) FROM operators`
	
	// Get total count
	var totalCount int
	err := r.db.Pool.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count operators: %w", err)
	}
	
	// Add pagination
	args := []interface{}{}
	if limit > 0 {
		query += " LIMIT $1"
		args = append(args, limit)
		if offset > 0 {
			query += " OFFSET $2"
			args = append(args, offset)
		}
	} else if offset > 0 {
		query += " OFFSET $1"
		args = append(args, offset)
	}
	
	// Execute query
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list operators: %w", err)
	}
	defer rows.Close()
	
	operators := []*models.Operator{}
	for rows.Next() {
		operator := &models.Operator{}
		err := rows.Scan(
			&operator.ID, &operator.Name, &operator.Code, &operator.ContactEmail,
			&operator.ContactPhone, &operator.Address, &operator.IsActive,
			&operator.Settings, &operator.CreatedAt, &operator.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan operator: %w", err)
		}
		operators = append(operators, operator)
	}
	
	return operators, totalCount, nil
}