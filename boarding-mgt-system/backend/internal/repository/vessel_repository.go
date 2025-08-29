package repository

import (
	"context"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type VesselRepository interface {
	Create(ctx context.Context, vessel *models.Vessel) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Vessel, error)
	GetByRegistration(ctx context.Context, regNumber string) (*models.Vessel, error)
	Update(ctx context.Context, vessel *models.Vessel) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByOperator(ctx context.Context, operatorID uuid.UUID, limit, offset int) ([]*models.Vessel, int, error)
	GetAvailableVessels(ctx context.Context, operatorID uuid.UUID) ([]*models.Vessel, error)
}

type vesselRepository struct {
	db *database.DB
}

func NewVesselRepository(db *database.DB) VesselRepository {
	return &vesselRepository{db: db}
}

func (r *vesselRepository) Create(ctx context.Context, vessel *models.Vessel) error {
	query := `
		INSERT INTO vessels (
			operator_id, name, registration_number, vessel_type,
			capacity, deck_count, seat_configuration, amenities
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, is_active, created_at, updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		vessel.OperatorID, vessel.Name, vessel.RegistrationNumber, vessel.VesselType,
		vessel.Capacity, vessel.DeckCount, vessel.SeatConfiguration, vessel.Amenities,
	).Scan(&vessel.ID, &vessel.IsActive, &vessel.CreatedAt, &vessel.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create vessel: %w", err)
	}
	
	return nil
}

func (r *vesselRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Vessel, error) {
	query := `
		SELECT 
			v.id, v.operator_id, v.name, v.registration_number, v.vessel_type,
			v.capacity, v.deck_count, v.seat_configuration, v.amenities,
			v.is_active, v.created_at, v.updated_at,
			o.id, o.name, o.code, o.contact_email, o.is_active
		FROM vessels v
		LEFT JOIN operators o ON v.operator_id = o.id
		WHERE v.id = $1
	`
	
	vessel := &models.Vessel{}
	operator := &models.Operator{}
	
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&vessel.ID, &vessel.OperatorID, &vessel.Name, &vessel.RegistrationNumber,
		&vessel.VesselType, &vessel.Capacity, &vessel.DeckCount,
		&vessel.SeatConfiguration, &vessel.Amenities,
		&vessel.IsActive, &vessel.CreatedAt, &vessel.UpdatedAt,
		&operator.ID, &operator.Name, &operator.Code, &operator.ContactEmail, &operator.IsActive,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("vessel not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get vessel: %w", err)
	}
	
	vessel.Operator = operator
	return vessel, nil
}

func (r *vesselRepository) GetByRegistration(ctx context.Context, regNumber string) (*models.Vessel, error) {
	query := `
		SELECT 
			id, operator_id, name, registration_number, vessel_type,
			capacity, deck_count, seat_configuration, amenities,
			is_active, created_at, updated_at
		FROM vessels
		WHERE registration_number = $1
	`
	
	vessel := &models.Vessel{}
	err := r.db.Pool.QueryRow(ctx, query, regNumber).Scan(
		&vessel.ID, &vessel.OperatorID, &vessel.Name, &vessel.RegistrationNumber,
		&vessel.VesselType, &vessel.Capacity, &vessel.DeckCount,
		&vessel.SeatConfiguration, &vessel.Amenities,
		&vessel.IsActive, &vessel.CreatedAt, &vessel.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("vessel not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get vessel: %w", err)
	}
	
	return vessel, nil
}

func (r *vesselRepository) Update(ctx context.Context, vessel *models.Vessel) error {
	query := `
		UPDATE vessels SET
			name = $2,
			vessel_type = $3,
			capacity = $4,
			deck_count = $5,
			seat_configuration = $6,
			amenities = $7,
			is_active = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		vessel.ID, vessel.Name, vessel.VesselType, vessel.Capacity,
		vessel.DeckCount, vessel.SeatConfiguration, vessel.Amenities, vessel.IsActive,
	).Scan(&vessel.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		return fmt.Errorf("vessel not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update vessel: %w", err)
	}
	
	return nil
}

func (r *vesselRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM vessels WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete vessel: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("vessel not found")
	}
	
	return nil
}

func (r *vesselRepository) ListByOperator(ctx context.Context, operatorID uuid.UUID, limit, offset int) ([]*models.Vessel, int, error) {
	query := `
		SELECT 
			id, operator_id, name, registration_number, vessel_type,
			capacity, deck_count, seat_configuration, amenities,
			is_active, created_at, updated_at
		FROM vessels
		WHERE operator_id = $1
		ORDER BY name ASC
	`
	countQuery := `SELECT COUNT(*) FROM vessels WHERE operator_id = $1`
	
	// Get total count
	var totalCount int
	err := r.db.Pool.QueryRow(ctx, countQuery, operatorID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vessels: %w", err)
	}
	
	// Add pagination
	args := []interface{}{operatorID}
	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
		if offset > 0 {
			query += " OFFSET $3"
			args = append(args, offset)
		}
	} else if offset > 0 {
		query += " OFFSET $2"
		args = append(args, offset)
	}
	
	// Execute query
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vessels: %w", err)
	}
	defer rows.Close()
	
	vessels := []*models.Vessel{}
	for rows.Next() {
		vessel := &models.Vessel{}
		err := rows.Scan(
			&vessel.ID, &vessel.OperatorID, &vessel.Name, &vessel.RegistrationNumber,
			&vessel.VesselType, &vessel.Capacity, &vessel.DeckCount,
			&vessel.SeatConfiguration, &vessel.Amenities,
			&vessel.IsActive, &vessel.CreatedAt, &vessel.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vessel: %w", err)
		}
		vessels = append(vessels, vessel)
	}
	
	return vessels, totalCount, nil
}

func (r *vesselRepository) GetAvailableVessels(ctx context.Context, operatorID uuid.UUID) ([]*models.Vessel, error) {
	query := `
		SELECT 
			id, operator_id, name, registration_number, vessel_type,
			capacity, deck_count, seat_configuration, amenities,
			is_active, created_at, updated_at
		FROM vessels
		WHERE operator_id = $1 AND is_active = true
		ORDER BY name ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, query, operatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available vessels: %w", err)
	}
	defer rows.Close()
	
	vessels := []*models.Vessel{}
	for rows.Next() {
		vessel := &models.Vessel{}
		err := rows.Scan(
			&vessel.ID, &vessel.OperatorID, &vessel.Name, &vessel.RegistrationNumber,
			&vessel.VesselType, &vessel.Capacity, &vessel.DeckCount,
			&vessel.SeatConfiguration, &vessel.Amenities,
			&vessel.IsActive, &vessel.CreatedAt, &vessel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vessel: %w", err)
		}
		vessels = append(vessels, vessel)
	}
	
	return vessels, nil
}