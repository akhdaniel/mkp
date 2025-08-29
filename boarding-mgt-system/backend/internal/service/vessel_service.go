package service

import (
	"context"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/ferryflow/boarding-mgt-system/internal/repository"
	"github.com/google/uuid"
)

type VesselService interface {
	CreateVessel(ctx context.Context, req *models.CreateVesselRequest) (*models.Vessel, error)
	GetVessel(ctx context.Context, id uuid.UUID) (*models.Vessel, error)
	UpdateVessel(ctx context.Context, id uuid.UUID, req *models.UpdateVesselRequest) (*models.Vessel, error)
	DeleteVessel(ctx context.Context, id uuid.UUID) error
	ListVesselsByOperator(ctx context.Context, operatorID uuid.UUID, limit, offset int) ([]*models.Vessel, int, error)
	GetAvailableVessels(ctx context.Context, operatorID uuid.UUID) ([]*models.Vessel, error)
}

type vesselService struct {
	vesselRepo   repository.VesselRepository
	operatorRepo repository.OperatorRepository
}

func NewVesselService(vesselRepo repository.VesselRepository, operatorRepo repository.OperatorRepository) VesselService {
	return &vesselService{
		vesselRepo:   vesselRepo,
		operatorRepo: operatorRepo,
	}
}

func (s *vesselService) CreateVessel(ctx context.Context, req *models.CreateVesselRequest) (*models.Vessel, error) {
	// Verify operator exists
	_, err := s.operatorRepo.GetByID(ctx, req.OperatorID)
	if err != nil {
		return nil, fmt.Errorf("operator not found: %w", err)
	}

	// Check if vessel with same registration exists
	existing, _ := s.vesselRepo.GetByRegistration(ctx, req.RegistrationNumber)
	if existing != nil {
		return nil, fmt.Errorf("vessel with registration %s already exists", req.RegistrationNumber)
	}

	vessel := &models.Vessel{
		OperatorID:         req.OperatorID,
		Name:               req.Name,
		RegistrationNumber: req.RegistrationNumber,
		VesselType:         req.VesselType,
		Capacity:           req.Capacity,
		DeckCount:          req.DeckCount,
		SeatConfiguration:  req.SeatConfiguration,
		IsActive:           true,
	}

	if req.Amenities != nil {
		vessel.Amenities = req.Amenities
	} else {
		vessel.Amenities = make(map[string]interface{})
	}

	if err := s.vesselRepo.Create(ctx, vessel); err != nil {
		return nil, fmt.Errorf("failed to create vessel: %w", err)
	}

	return vessel, nil
}

func (s *vesselService) GetVessel(ctx context.Context, id uuid.UUID) (*models.Vessel, error) {
	vessel, err := s.vesselRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("vessel not found: %w", err)
	}

	return vessel, nil
}

func (s *vesselService) UpdateVessel(ctx context.Context, id uuid.UUID, req *models.UpdateVesselRequest) (*models.Vessel, error) {
	// Get existing vessel
	vessel, err := s.vesselRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("vessel not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		vessel.Name = *req.Name
	}
	if req.VesselType != nil {
		vessel.VesselType = *req.VesselType
	}
	if req.Capacity != nil {
		vessel.Capacity = *req.Capacity
	}
	if req.DeckCount != nil {
		vessel.DeckCount = *req.DeckCount
	}
	if req.SeatConfiguration != nil {
		vessel.SeatConfiguration = req.SeatConfiguration
	}
	if req.Amenities != nil {
		vessel.Amenities = req.Amenities
	}
	if req.IsActive != nil {
		vessel.IsActive = *req.IsActive
	}

	if err := s.vesselRepo.Update(ctx, vessel); err != nil {
		return nil, fmt.Errorf("failed to update vessel: %w", err)
	}

	return vessel, nil
}

func (s *vesselService) DeleteVessel(ctx context.Context, id uuid.UUID) error {
	if err := s.vesselRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete vessel: %w", err)
	}

	return nil
}

func (s *vesselService) ListVesselsByOperator(ctx context.Context, operatorID uuid.UUID, limit, offset int) ([]*models.Vessel, int, error) {
	vessels, total, err := s.vesselRepo.ListByOperator(ctx, operatorID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list vessels: %w", err)
	}

	return vessels, total, nil
}

func (s *vesselService) GetAvailableVessels(ctx context.Context, operatorID uuid.UUID) ([]*models.Vessel, error) {
	vessels, err := s.vesselRepo.GetAvailableVessels(ctx, operatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available vessels: %w", err)
	}

	return vessels, nil
}