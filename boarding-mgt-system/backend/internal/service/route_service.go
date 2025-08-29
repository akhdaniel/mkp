package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/ferryflow/boarding-mgt-system/internal/repository"
	"github.com/google/uuid"
)

type RouteService interface {
	CreateRoute(ctx context.Context, req *models.CreateRouteRequest) (*models.Route, error)
	GetRoute(ctx context.Context, id uuid.UUID) (*models.Route, error)
	UpdateRoute(ctx context.Context, id uuid.UUID, req *models.UpdateRouteRequest) (*models.Route, error)
	DeleteRoute(ctx context.Context, id uuid.UUID) error
	ListRoutesByOperator(ctx context.Context, operatorID uuid.UUID) ([]*models.Route, error)
	SearchRoutes(ctx context.Context, departurePortID, arrivalPortID uuid.UUID) ([]*models.Route, error)
}

type routeService struct {
	routeRepo repository.RouteRepository
	portRepo  repository.PortRepository
}

func NewRouteService(routeRepo repository.RouteRepository, portRepo repository.PortRepository) RouteService {
	return &routeService{
		routeRepo: routeRepo,
		portRepo:  portRepo,
	}
}

func (s *routeService) CreateRoute(ctx context.Context, req *models.CreateRouteRequest) (*models.Route, error) {
	// Verify ports exist
	_, err := s.portRepo.GetByID(ctx, req.DeparturePortID)
	if err != nil {
		return nil, fmt.Errorf("departure port not found: %w", err)
	}

	_, err = s.portRepo.GetByID(ctx, req.ArrivalPortID)
	if err != nil {
		return nil, fmt.Errorf("arrival port not found: %w", err)
	}

	// Check if departure and arrival ports are different
	if req.DeparturePortID == req.ArrivalPortID {
		return nil, fmt.Errorf("departure and arrival ports must be different")
	}

	// Parse duration
	duration, err := time.ParseDuration(req.EstimatedDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid duration format: %w", err)
	}

	route := &models.Route{
		OperatorID:        req.OperatorID,
		Name:              req.Name,
		DeparturePortID:   req.DeparturePortID,
		ArrivalPortID:     req.ArrivalPortID,
		EstimatedDuration: duration,
		IsActive:          true,
	}

	if req.DistanceKM > 0 {
		route.DistanceKM = &req.DistanceKM
	}

	if err := s.routeRepo.Create(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	return route, nil
}

func (s *routeService) GetRoute(ctx context.Context, id uuid.UUID) (*models.Route, error) {
	route, err := s.routeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}

	return route, nil
}

func (s *routeService) UpdateRoute(ctx context.Context, id uuid.UUID, req *models.UpdateRouteRequest) (*models.Route, error) {
	// Get existing route
	route, err := s.routeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		route.Name = *req.Name
	}
	if req.DistanceKM != nil {
		route.DistanceKM = req.DistanceKM
	}
	if req.EstimatedDuration != nil {
		duration, err := time.ParseDuration(*req.EstimatedDuration)
		if err != nil {
			return nil, fmt.Errorf("invalid duration format: %w", err)
		}
		route.EstimatedDuration = duration
	}
	if req.IsActive != nil {
		route.IsActive = *req.IsActive
	}

	if err := s.routeRepo.Update(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to update route: %w", err)
	}

	return route, nil
}

func (s *routeService) DeleteRoute(ctx context.Context, id uuid.UUID) error {
	if err := s.routeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete route: %w", err)
	}

	return nil
}

func (s *routeService) ListRoutesByOperator(ctx context.Context, operatorID uuid.UUID) ([]*models.Route, error) {
	routes, err := s.routeRepo.ListByOperator(ctx, operatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}

	return routes, nil
}

func (s *routeService) SearchRoutes(ctx context.Context, departurePortID, arrivalPortID uuid.UUID) ([]*models.Route, error) {
	routes, err := s.routeRepo.SearchRoutes(ctx, departurePortID, arrivalPortID)
	if err != nil {
		return nil, fmt.Errorf("failed to search routes: %w", err)
	}

	return routes, nil
}