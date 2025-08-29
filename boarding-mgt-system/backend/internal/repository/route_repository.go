package repository

import (
	"context"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RouteRepository interface {
	Create(ctx context.Context, route *models.Route) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Route, error)
	Update(ctx context.Context, route *models.Route) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByOperator(ctx context.Context, operatorID uuid.UUID) ([]*models.Route, error)
	SearchRoutes(ctx context.Context, departurePortID, arrivalPortID uuid.UUID) ([]*models.Route, error)
}

type routeRepository struct {
	db *database.DB
}

func NewRouteRepository(db *database.DB) RouteRepository {
	return &routeRepository{db: db}
}

func (r *routeRepository) Create(ctx context.Context, route *models.Route) error {
	query := `
		INSERT INTO routes (
			operator_id, name, departure_port_id, arrival_port_id,
			distance_km, estimated_duration
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, is_active, created_at, updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		route.OperatorID, route.Name, route.DeparturePortID, route.ArrivalPortID,
		route.DistanceKM, route.EstimatedDuration,
	).Scan(&route.ID, &route.IsActive, &route.CreatedAt, &route.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create route: %w", err)
	}
	
	return nil
}

func (r *routeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Route, error) {
	query := `
		SELECT 
			r.id, r.operator_id, r.name, r.departure_port_id, r.arrival_port_id,
			r.distance_km, r.estimated_duration, r.is_active, r.created_at, r.updated_at,
			o.id, o.name, o.code,
			dp.id, dp.name, dp.code, dp.city, dp.country,
			ap.id, ap.name, ap.code, ap.city, ap.country
		FROM routes r
		LEFT JOIN operators o ON r.operator_id = o.id
		LEFT JOIN ports dp ON r.departure_port_id = dp.id
		LEFT JOIN ports ap ON r.arrival_port_id = ap.id
		WHERE r.id = $1
	`
	
	route := &models.Route{}
	operator := &models.Operator{}
	departurePort := &models.Port{}
	arrivalPort := &models.Port{}
	
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&route.ID, &route.OperatorID, &route.Name, &route.DeparturePortID, &route.ArrivalPortID,
		&route.DistanceKM, &route.EstimatedDuration, &route.IsActive, &route.CreatedAt, &route.UpdatedAt,
		&operator.ID, &operator.Name, &operator.Code,
		&departurePort.ID, &departurePort.Name, &departurePort.Code, &departurePort.City, &departurePort.Country,
		&arrivalPort.ID, &arrivalPort.Name, &arrivalPort.Code, &arrivalPort.City, &arrivalPort.Country,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("route not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get route: %w", err)
	}
	
	route.Operator = operator
	route.DeparturePort = departurePort
	route.ArrivalPort = arrivalPort
	
	return route, nil
}

func (r *routeRepository) Update(ctx context.Context, route *models.Route) error {
	query := `
		UPDATE routes SET
			name = $2,
			distance_km = $3,
			estimated_duration = $4,
			is_active = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		route.ID, route.Name, route.DistanceKM,
		route.EstimatedDuration, route.IsActive,
	).Scan(&route.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		return fmt.Errorf("route not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}
	
	return nil
}

func (r *routeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM routes WHERE id = $1`
	
	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("route not found")
	}
	
	return nil
}

func (r *routeRepository) ListByOperator(ctx context.Context, operatorID uuid.UUID) ([]*models.Route, error) {
	query := `
		SELECT 
			r.id, r.operator_id, r.name, r.departure_port_id, r.arrival_port_id,
			r.distance_km, r.estimated_duration, r.is_active, r.created_at, r.updated_at,
			dp.id, dp.name, dp.code, dp.city, dp.country,
			ap.id, ap.name, ap.code, ap.city, ap.country
		FROM routes r
		LEFT JOIN ports dp ON r.departure_port_id = dp.id
		LEFT JOIN ports ap ON r.arrival_port_id = ap.id
		WHERE r.operator_id = $1 AND r.is_active = true
		ORDER BY r.name ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, query, operatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}
	defer rows.Close()
	
	routes := []*models.Route{}
	for rows.Next() {
		route := &models.Route{}
		departurePort := &models.Port{}
		arrivalPort := &models.Port{}
		
		err := rows.Scan(
			&route.ID, &route.OperatorID, &route.Name, &route.DeparturePortID, &route.ArrivalPortID,
			&route.DistanceKM, &route.EstimatedDuration, &route.IsActive, &route.CreatedAt, &route.UpdatedAt,
			&departurePort.ID, &departurePort.Name, &departurePort.Code, &departurePort.City, &departurePort.Country,
			&arrivalPort.ID, &arrivalPort.Name, &arrivalPort.Code, &arrivalPort.City, &arrivalPort.Country,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan route: %w", err)
		}
		
		route.DeparturePort = departurePort
		route.ArrivalPort = arrivalPort
		routes = append(routes, route)
	}
	
	return routes, nil
}

func (r *routeRepository) SearchRoutes(ctx context.Context, departurePortID, arrivalPortID uuid.UUID) ([]*models.Route, error) {
	query := `
		SELECT 
			r.id, r.operator_id, r.name, r.departure_port_id, r.arrival_port_id,
			r.distance_km, r.estimated_duration, r.is_active, r.created_at, r.updated_at,
			o.id, o.name, o.code,
			dp.id, dp.name, dp.code, dp.city, dp.country,
			ap.id, ap.name, ap.code, ap.city, ap.country
		FROM routes r
		LEFT JOIN operators o ON r.operator_id = o.id
		LEFT JOIN ports dp ON r.departure_port_id = dp.id
		LEFT JOIN ports ap ON r.arrival_port_id = ap.id
		WHERE r.departure_port_id = $1 AND r.arrival_port_id = $2 AND r.is_active = true
		ORDER BY r.name ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, query, departurePortID, arrivalPortID)
	if err != nil {
		return nil, fmt.Errorf("failed to search routes: %w", err)
	}
	defer rows.Close()
	
	routes := []*models.Route{}
	for rows.Next() {
		route := &models.Route{}
		operator := &models.Operator{}
		departurePort := &models.Port{}
		arrivalPort := &models.Port{}
		
		err := rows.Scan(
			&route.ID, &route.OperatorID, &route.Name, &route.DeparturePortID, &route.ArrivalPortID,
			&route.DistanceKM, &route.EstimatedDuration, &route.IsActive, &route.CreatedAt, &route.UpdatedAt,
			&operator.ID, &operator.Name, &operator.Code,
			&departurePort.ID, &departurePort.Name, &departurePort.Code, &departurePort.City, &departurePort.Country,
			&arrivalPort.ID, &arrivalPort.Name, &arrivalPort.Code, &arrivalPort.City, &arrivalPort.Country,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan route: %w", err)
		}
		
		route.Operator = operator
		route.DeparturePort = departurePort
		route.ArrivalPort = arrivalPort
		routes = append(routes, route)
	}
	
	return routes, nil
}