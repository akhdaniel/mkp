package service

import (
	"context"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/ferryflow/boarding-mgt-system/internal/repository"
	"github.com/google/uuid"
)

type OperatorService interface {
	CreateOperator(ctx context.Context, req *models.CreateOperatorRequest) (*models.Operator, error)
	GetOperator(ctx context.Context, id uuid.UUID) (*models.Operator, error)
	UpdateOperator(ctx context.Context, id uuid.UUID, req *models.UpdateOperatorRequest) (*models.Operator, error)
	DeleteOperator(ctx context.Context, id uuid.UUID) error
	ListOperators(ctx context.Context, limit, offset int) ([]*models.Operator, int, error)
}

type operatorService struct {
	operatorRepo repository.OperatorRepository
}

func NewOperatorService(operatorRepo repository.OperatorRepository) OperatorService {
	return &operatorService{
		operatorRepo: operatorRepo,
	}
}

func (s *operatorService) CreateOperator(ctx context.Context, req *models.CreateOperatorRequest) (*models.Operator, error) {
	// Check if operator with same code exists
	existing, _ := s.operatorRepo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("operator with code %s already exists", req.Code)
	}

	operator := &models.Operator{
		Name:         req.Name,
		Code:         req.Code,
		ContactEmail: req.ContactEmail,
		IsActive:     true,
	}

	if req.ContactPhone != "" {
		operator.ContactPhone = &req.ContactPhone
	}
	if req.Address != "" {
		operator.Address = &req.Address
	}
	if req.Settings != nil {
		operator.Settings = req.Settings
	} else {
		operator.Settings = make(map[string]interface{})
	}

	if err := s.operatorRepo.Create(ctx, operator); err != nil {
		return nil, fmt.Errorf("failed to create operator: %w", err)
	}

	return operator, nil
}

func (s *operatorService) GetOperator(ctx context.Context, id uuid.UUID) (*models.Operator, error) {
	operator, err := s.operatorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("operator not found: %w", err)
	}

	return operator, nil
}

func (s *operatorService) UpdateOperator(ctx context.Context, id uuid.UUID, req *models.UpdateOperatorRequest) (*models.Operator, error) {
	// Get existing operator
	operator, err := s.operatorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("operator not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		operator.Name = *req.Name
	}
	if req.ContactEmail != nil {
		operator.ContactEmail = *req.ContactEmail
	}
	if req.ContactPhone != nil {
		operator.ContactPhone = req.ContactPhone
	}
	if req.Address != nil {
		operator.Address = req.Address
	}
	if req.IsActive != nil {
		operator.IsActive = *req.IsActive
	}
	if req.Settings != nil {
		operator.Settings = req.Settings
	}

	if err := s.operatorRepo.Update(ctx, operator); err != nil {
		return nil, fmt.Errorf("failed to update operator: %w", err)
	}

	return operator, nil
}

func (s *operatorService) DeleteOperator(ctx context.Context, id uuid.UUID) error {
	if err := s.operatorRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete operator: %w", err)
	}

	return nil
}

func (s *operatorService) ListOperators(ctx context.Context, limit, offset int) ([]*models.Operator, int, error) {
	operators, total, err := s.operatorRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list operators: %w", err)
	}

	return operators, total, nil
}