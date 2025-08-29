package repository

import (
	"context"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PortRepository interface {
	Create(ctx context.Context, port *models.Port) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Port, error)
	GetByCode(ctx context.Context, code string) (*models.Port, error)
	Update(ctx context.Context, port *models.Port) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*models.Port, int, error)
	SearchByLocation(ctx context.Context, city, country string) ([]*models.Port, error)
}

type portRepository struct {
	db *database.DB
}

func NewPortRepository(db *database.DB) PortRepository {
	return &portRepository{db: db}
}

func (r *portRepository) Create(ctx context.Context, port *models.Port) error {
	query := `
		INSERT INTO ports (
			name, code, city, country, timezone, coordinates, facilities
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, is_active, created_at, updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		port.Name, port.Code, port.City, port.Country,
		port.Timezone, port.Coordinates, port.Facilities,
	).Scan(&port.ID, &port.IsActive, &port.CreatedAt, &port.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create port: %w", err)
	}
	
	return nil
}

func (r *portRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Port, error) {
	query := `
		SELECT 
			id, name, code, city, country, timezone,
			coordinates, facilities, is_active, created_at, updated_at
		FROM ports
		WHERE id = $1
	`
	
	port := &models.Port{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&port.ID, &port.Name, &port.Code, &port.City, &port.Country,
		&port.Timezone, &port.Coordinates, &port.Facilities,
		&port.IsActive, &port.CreatedAt, &port.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("port not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get port: %w", err)
	}
	
	return port, nil
}

func (r *portRepository) GetByCode(ctx context.Context, code string) (*models.Port, error) {
	query := `
		SELECT 
			id, name, code, city, country, timezone,
			coordinates, facilities, is_active, created_at, updated_at
		FROM ports
		WHERE code = $1
	`
	
	port := &models.Port{}
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(
		&port.ID, &port.Name, &port.Code, &port.City, &port.Country,
		&port.Timezone, &port.Coordinates, &port.Facilities,
		&port.IsActive, &port.CreatedAt, &port.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("port not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get port: %w", err)
	}
	
	return port, nil
}

func (r *portRepository) Update(ctx context.Context, port *models.Port) error {
	query := `
		UPDATE ports SET
			name = $2,
			city = $3,
			country = $4,
			timezone = $5,
			coordinates = $6,
			facilities = $7,
			is_active = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		port.ID, port.Name, port.City, port.Country,
		port.Timezone, port.Coordinates, port.Facilities, port.IsActive,
	).Scan(&port.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		return fmt.Errorf("port not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update port: %w", err)
	}
	
	return nil
}

func (r *portRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM ports WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete port: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("port not found")
	}
	
	return nil
}

func (r *portRepository) List(ctx context.Context, limit, offset int) ([]*models.Port, int, error) {
	query := `
		SELECT 
			id, name, code, city, country, timezone,
			coordinates, facilities, is_active, created_at, updated_at
		FROM ports
		WHERE is_active = true
		ORDER BY name ASC
	`
	countQuery := `SELECT COUNT(*) FROM ports WHERE is_active = true`
	
	// Get total count
	var totalCount int
	err := r.db.Pool.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count ports: %w", err)
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
		return nil, 0, fmt.Errorf("failed to list ports: %w", err)
	}
	defer rows.Close()
	
	ports := []*models.Port{}
	for rows.Next() {
		port := &models.Port{}
		err := rows.Scan(
			&port.ID, &port.Name, &port.Code, &port.City, &port.Country,
			&port.Timezone, &port.Coordinates, &port.Facilities,
			&port.IsActive, &port.CreatedAt, &port.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan port: %w", err)
		}
		ports = append(ports, port)
	}
	
	return ports, totalCount, nil
}

func (r *portRepository) SearchByLocation(ctx context.Context, city, country string) ([]*models.Port, error) {
	query := `
		SELECT 
			id, name, code, city, country, timezone,
			coordinates, facilities, is_active, created_at, updated_at
		FROM ports
		WHERE is_active = true
	`
	
	args := []interface{}{}
	argCount := 0
	
	if city != "" {
		argCount++
		query += fmt.Sprintf(" AND city ILIKE $%d", argCount)
		args = append(args, "%"+city+"%")
	}
	
	if country != "" {
		argCount++
		query += fmt.Sprintf(" AND country ILIKE $%d", argCount)
		args = append(args, "%"+country+"%")
	}
	
	query += " ORDER BY name ASC"
	
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search ports: %w", err)
	}
	defer rows.Close()
	
	ports := []*models.Port{}
	for rows.Next() {
		port := &models.Port{}
		err := rows.Scan(
			&port.ID, &port.Name, &port.Code, &port.City, &port.Country,
			&port.Timezone, &port.Coordinates, &port.Facilities,
			&port.IsActive, &port.CreatedAt, &port.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan port: %w", err)
		}
		ports = append(ports, port)
	}
	
	return ports, nil
}