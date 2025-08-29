package service

import (
	"context"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/ferryflow/boarding-mgt-system/internal/repository"
	"github.com/google/uuid"
)

type PortService interface {
	CreatePort(ctx context.Context, req *models.CreatePortRequest) (*models.Port, error)
	GetPort(ctx context.Context, id uuid.UUID) (*models.Port, error)
	UpdatePort(ctx context.Context, id uuid.UUID, req *models.UpdatePortRequest) (*models.Port, error)
	DeletePort(ctx context.Context, id uuid.UUID) error
	ListPorts(ctx context.Context, limit, offset int) ([]*models.Port, int, error)
	SearchPorts(ctx context.Context, city, country string) ([]*models.Port, error)
}

type portService struct {
	portRepo repository.PortRepository
}

func NewPortService(portRepo repository.PortRepository) PortService {
	return &portService{
		portRepo: portRepo,
	}
}

func (s *portService) CreatePort(ctx context.Context, req *models.CreatePortRequest) (*models.Port, error) {
	// Check if port with same code exists
	existing, _ := s.portRepo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("port with code %s already exists", req.Code)
	}

	port := &models.Port{
		Name:        req.Name,
		Code:        req.Code,
		City:        req.City,
		Country:     req.Country,
		Timezone:    req.Timezone,
		Coordinates: req.Coordinates,
		IsActive:    true,
	}

	if req.Facilities != nil {
		port.Facilities = req.Facilities
	} else {
		port.Facilities = make(map[string]interface{})
	}

	if err := s.portRepo.Create(ctx, port); err != nil {
		return nil, fmt.Errorf("failed to create port: %w", err)
	}

	return port, nil
}

func (s *portService) GetPort(ctx context.Context, id uuid.UUID) (*models.Port, error) {
	port, err := s.portRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("port not found: %w", err)
	}

	return port, nil
}

func (s *portService) UpdatePort(ctx context.Context, id uuid.UUID, req *models.UpdatePortRequest) (*models.Port, error) {
	// Get existing port
	port, err := s.portRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("port not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		port.Name = *req.Name
	}
	if req.City != nil {
		port.City = *req.City
	}
	if req.Country != nil {
		port.Country = *req.Country
	}
	if req.Timezone != nil {
		port.Timezone = *req.Timezone
	}
	if req.Coordinates != nil {
		port.Coordinates = req.Coordinates
	}
	if req.Facilities != nil {
		port.Facilities = req.Facilities
	}
	if req.IsActive != nil {
		port.IsActive = *req.IsActive
	}

	if err := s.portRepo.Update(ctx, port); err != nil {
		return nil, fmt.Errorf("failed to update port: %w", err)
	}

	return port, nil
}

func (s *portService) DeletePort(ctx context.Context, id uuid.UUID) error {
	if err := s.portRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete port: %w", err)
	}

	return nil
}

func (s *portService) ListPorts(ctx context.Context, limit, offset int) ([]*models.Port, int, error) {
	ports, total, err := s.portRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list ports: %w", err)
	}

	return ports, total, nil
}

func (s *portService) SearchPorts(ctx context.Context, city, country string) ([]*models.Port, error) {
	ports, err := s.portRepo.SearchByLocation(ctx, city, country)
	if err != nil {
		return nil, fmt.Errorf("failed to search ports: %w", err)
	}

	return ports, nil
}