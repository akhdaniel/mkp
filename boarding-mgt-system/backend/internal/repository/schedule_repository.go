package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *models.Schedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Schedule, error)
	Update(ctx context.Context, schedule *models.Schedule) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, reason *string) error
	Delete(ctx context.Context, id uuid.UUID) error
	SearchSchedules(ctx context.Context, req *models.SearchScheduleRequest) ([]*models.Schedule, int, error)
	GetByOperatorAndDate(ctx context.Context, operatorID uuid.UUID, date time.Time) ([]*models.Schedule, error)
	GetUpcomingSchedules(ctx context.Context, limit int) ([]*models.Schedule, error)
}

type scheduleRepository struct {
	db *database.DB
}

func NewScheduleRepository(db *database.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(ctx context.Context, schedule *models.Schedule) error {
	query := `
		INSERT INTO schedules (
			operator_id, route_id, vessel_id, departure_date, departure_time,
			arrival_time, base_price, total_capacity, available_seats
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, status, version, created_at, updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		schedule.OperatorID, schedule.RouteID, schedule.VesselID,
		schedule.DepartureDate, schedule.DepartureTime, schedule.ArrivalTime,
		schedule.BasePrice, schedule.TotalCapacity, schedule.AvailableSeats,
	).Scan(&schedule.ID, &schedule.Status, &schedule.Version, &schedule.CreatedAt, &schedule.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}
	
	return nil
}

func (r *scheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Schedule, error) {
	query := `
		SELECT 
			s.id, s.operator_id, s.route_id, s.vessel_id, s.departure_date,
			s.departure_time, s.arrival_time, s.base_price, s.total_capacity,
			s.available_seats, s.status, s.cancellation_reason, s.version,
			s.created_at, s.updated_at,
			o.id, o.name, o.code,
			r.id, r.name, r.departure_port_id, r.arrival_port_id,
			v.id, v.name, v.registration_number, v.capacity
		FROM schedules s
		LEFT JOIN operators o ON s.operator_id = o.id
		LEFT JOIN routes r ON s.route_id = r.id
		LEFT JOIN vessels v ON s.vessel_id = v.id
		WHERE s.id = $1
	`
	
	schedule := &models.Schedule{}
	operator := &models.Operator{}
	route := &models.Route{}
	vessel := &models.Vessel{}
	
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&schedule.ID, &schedule.OperatorID, &schedule.RouteID, &schedule.VesselID,
		&schedule.DepartureDate, &schedule.DepartureTime, &schedule.ArrivalTime,
		&schedule.BasePrice, &schedule.TotalCapacity, &schedule.AvailableSeats,
		&schedule.Status, &schedule.CancellationReason, &schedule.Version,
		&schedule.CreatedAt, &schedule.UpdatedAt,
		&operator.ID, &operator.Name, &operator.Code,
		&route.ID, &route.Name, &route.DeparturePortID, &route.ArrivalPortID,
		&vessel.ID, &vessel.Name, &vessel.RegistrationNumber, &vessel.Capacity,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("schedule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}
	
	schedule.Operator = operator
	schedule.Route = route
	schedule.Vessel = vessel
	
	return schedule, nil
}

func (r *scheduleRepository) Update(ctx context.Context, schedule *models.Schedule) error {
	query := `
		UPDATE schedules SET
			departure_date = $2,
			departure_time = $3,
			arrival_time = $4,
			base_price = $5,
			status = $6,
			version = version + 1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND version = $7
		RETURNING version, updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		schedule.ID, schedule.DepartureDate, schedule.DepartureTime,
		schedule.ArrivalTime, schedule.BasePrice, schedule.Status, schedule.Version,
	).Scan(&schedule.Version, &schedule.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		return fmt.Errorf("schedule not found or version mismatch")
	}
	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}
	
	return nil
}

func (r *scheduleRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, reason *string) error {
	query := `
		UPDATE schedules SET
			status = $2,
			cancellation_reason = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err := r.db.Pool.Exec(ctx, query, id, status, reason)
	if err != nil {
		return fmt.Errorf("failed to update schedule status: %w", err)
	}
	
	return nil
}

func (r *scheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedules WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("schedule not found")
	}
	
	return nil
}

func (r *scheduleRepository) SearchSchedules(ctx context.Context, req *models.SearchScheduleRequest) ([]*models.Schedule, int, error) {
	query := `
		SELECT 
			s.id, s.operator_id, s.route_id, s.vessel_id, s.departure_date,
			s.departure_time, s.arrival_time, s.base_price, s.total_capacity,
			s.available_seats, s.status, s.cancellation_reason, s.version,
			s.created_at, s.updated_at
		FROM schedules s
		JOIN routes r ON s.route_id = r.id
		WHERE r.departure_port_id = $1 
			AND r.arrival_port_id = $2
			AND s.departure_date = $3
			AND s.status = 'scheduled'
			AND s.available_seats >= $4
		ORDER BY s.departure_time ASC
	`
	
	passengerCount := req.PassengerCount
	if passengerCount == 0 {
		passengerCount = 1
	}
	
	args := []interface{}{
		req.DeparturePortID,
		req.ArrivalPortID,
		req.DepartureDate,
		passengerCount,
	}
	
	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM schedules s
		JOIN routes r ON s.route_id = r.id
		WHERE r.departure_port_id = $1 
			AND r.arrival_port_id = $2
			AND s.departure_date = $3
			AND s.status = 'scheduled'
			AND s.available_seats >= $4
	`
	
	var totalCount int
	err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count schedules: %w", err)
	}
	
	// Add pagination
	if req.Limit > 0 {
		query += " LIMIT $5"
		args = append(args, req.Limit)
		if req.Offset > 0 {
			query += " OFFSET $6"
			args = append(args, req.Offset)
		}
	} else if req.Offset > 0 {
		query += " OFFSET $5"
		args = append(args, req.Offset)
	}
	
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search schedules: %w", err)
	}
	defer rows.Close()
	
	schedules := []*models.Schedule{}
	for rows.Next() {
		schedule := &models.Schedule{}
		err := rows.Scan(
			&schedule.ID, &schedule.OperatorID, &schedule.RouteID, &schedule.VesselID,
			&schedule.DepartureDate, &schedule.DepartureTime, &schedule.ArrivalTime,
			&schedule.BasePrice, &schedule.TotalCapacity, &schedule.AvailableSeats,
			&schedule.Status, &schedule.CancellationReason, &schedule.Version,
			&schedule.CreatedAt, &schedule.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}
	
	return schedules, totalCount, nil
}

func (r *scheduleRepository) GetByOperatorAndDate(ctx context.Context, operatorID uuid.UUID, date time.Time) ([]*models.Schedule, error) {
	query := `
		SELECT 
			id, operator_id, route_id, vessel_id, departure_date,
			departure_time, arrival_time, base_price, total_capacity,
			available_seats, status, cancellation_reason, version,
			created_at, updated_at
		FROM schedules
		WHERE operator_id = $1 AND departure_date = $2
		ORDER BY departure_time ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, query, operatorID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules: %w", err)
	}
	defer rows.Close()
	
	schedules := []*models.Schedule{}
	for rows.Next() {
		schedule := &models.Schedule{}
		err := rows.Scan(
			&schedule.ID, &schedule.OperatorID, &schedule.RouteID, &schedule.VesselID,
			&schedule.DepartureDate, &schedule.DepartureTime, &schedule.ArrivalTime,
			&schedule.BasePrice, &schedule.TotalCapacity, &schedule.AvailableSeats,
			&schedule.Status, &schedule.CancellationReason, &schedule.Version,
			&schedule.CreatedAt, &schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}
	
	return schedules, nil
}

func (r *scheduleRepository) GetUpcomingSchedules(ctx context.Context, limit int) ([]*models.Schedule, error) {
	query := `
		SELECT 
			id, operator_id, route_id, vessel_id, departure_date,
			departure_time, arrival_time, base_price, total_capacity,
			available_seats, status, cancellation_reason, version,
			created_at, updated_at
		FROM schedules
		WHERE departure_date >= CURRENT_DATE
			AND status = 'scheduled'
		ORDER BY departure_date ASC, departure_time ASC
		LIMIT $1
	`
	
	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming schedules: %w", err)
	}
	defer rows.Close()
	
	schedules := []*models.Schedule{}
	for rows.Next() {
		schedule := &models.Schedule{}
		err := rows.Scan(
			&schedule.ID, &schedule.OperatorID, &schedule.RouteID, &schedule.VesselID,
			&schedule.DepartureDate, &schedule.DepartureTime, &schedule.ArrivalTime,
			&schedule.BasePrice, &schedule.TotalCapacity, &schedule.AvailableSeats,
			&schedule.Status, &schedule.CancellationReason, &schedule.Version,
			&schedule.CreatedAt, &schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}
	
	return schedules, nil
}