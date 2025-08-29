package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/ferryflow/boarding-mgt-system/internal/repository"
	"github.com/google/uuid"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, req *models.CreateScheduleRequest) (*models.Schedule, error)
	GetSchedule(ctx context.Context, id uuid.UUID) (*models.Schedule, error)
	UpdateSchedule(ctx context.Context, id uuid.UUID, req *models.UpdateScheduleRequest) (*models.Schedule, error)
	CancelSchedule(ctx context.Context, id uuid.UUID, reason string) error
	SearchSchedules(ctx context.Context, req *models.SearchScheduleRequest) ([]*models.Schedule, int, error)
	GetOperatorSchedules(ctx context.Context, operatorID uuid.UUID, date time.Time) ([]*models.Schedule, error)
	GetUpcomingSchedules(ctx context.Context, limit int) ([]*models.Schedule, error)
}

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
	routeRepo    repository.RouteRepository
	vesselRepo   repository.VesselRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, routeRepo repository.RouteRepository, vesselRepo repository.VesselRepository) ScheduleService {
	return &scheduleService{
		scheduleRepo: scheduleRepo,
		routeRepo:    routeRepo,
		vesselRepo:   vesselRepo,
	}
}

func (s *scheduleService) CreateSchedule(ctx context.Context, req *models.CreateScheduleRequest) (*models.Schedule, error) {
	// Verify route exists
	route, err := s.routeRepo.GetByID(ctx, req.RouteID)
	if err != nil {
		return nil, fmt.Errorf("route not found: %w", err)
	}

	// Verify vessel exists and belongs to operator
	vessel, err := s.vesselRepo.GetByID(ctx, req.VesselID)
	if err != nil {
		return nil, fmt.Errorf("vessel not found: %w", err)
	}

	if vessel.OperatorID != req.OperatorID {
		return nil, fmt.Errorf("vessel does not belong to operator")
	}

	// Parse dates and times
	departureDate, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		return nil, fmt.Errorf("invalid departure date format: %w", err)
	}

	departureTime, err := time.Parse("15:04", req.DepartureTime)
	if err != nil {
		return nil, fmt.Errorf("invalid departure time format: %w", err)
	}

	arrivalTime, err := time.Parse("15:04", req.ArrivalTime)
	if err != nil {
		return nil, fmt.Errorf("invalid arrival time format: %w", err)
	}

	schedule := &models.Schedule{
		OperatorID:     req.OperatorID,
		RouteID:        req.RouteID,
		VesselID:       req.VesselID,
		DepartureDate:  departureDate,
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		BasePrice:      req.BasePrice,
		TotalCapacity:  vessel.Capacity,
		AvailableSeats: vessel.Capacity,
		Status:         "scheduled",
	}

	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	// Attach related data
	schedule.Route = route
	schedule.Vessel = vessel

	return schedule, nil
}

func (s *scheduleService) GetSchedule(ctx context.Context, id uuid.UUID) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("schedule not found: %w", err)
	}

	return schedule, nil
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, id uuid.UUID, req *models.UpdateScheduleRequest) (*models.Schedule, error) {
	// Get existing schedule
	schedule, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("schedule not found: %w", err)
	}

	// Check if schedule can be modified
	if schedule.Status == "departed" || schedule.Status == "completed" {
		return nil, fmt.Errorf("cannot modify schedule with status %s", schedule.Status)
	}

	// Update fields if provided
	if req.DepartureDate != nil {
		departureDate, err := time.Parse("2006-01-02", *req.DepartureDate)
		if err != nil {
			return nil, fmt.Errorf("invalid departure date format: %w", err)
		}
		schedule.DepartureDate = departureDate
	}

	if req.DepartureTime != nil {
		departureTime, err := time.Parse("15:04", *req.DepartureTime)
		if err != nil {
			return nil, fmt.Errorf("invalid departure time format: %w", err)
		}
		schedule.DepartureTime = departureTime
	}

	if req.ArrivalTime != nil {
		arrivalTime, err := time.Parse("15:04", *req.ArrivalTime)
		if err != nil {
			return nil, fmt.Errorf("invalid arrival time format: %w", err)
		}
		schedule.ArrivalTime = arrivalTime
	}

	if req.BasePrice != nil {
		schedule.BasePrice = *req.BasePrice
	}

	if req.Status != nil {
		schedule.Status = *req.Status
	}

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	return schedule, nil
}

func (s *scheduleService) CancelSchedule(ctx context.Context, id uuid.UUID, reason string) error {
	if err := s.scheduleRepo.UpdateStatus(ctx, id, "cancelled", &reason); err != nil {
		return fmt.Errorf("failed to cancel schedule: %w", err)
	}

	// TODO: Notify affected passengers and process refunds

	return nil
}

func (s *scheduleService) SearchSchedules(ctx context.Context, req *models.SearchScheduleRequest) ([]*models.Schedule, int, error) {
	// Parse departure date
	departureDate, err := time.Parse("2006-01-02", req.DepartureDate)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid departure date format: %w", err)
	}

	searchReq := &models.SearchScheduleRequest{
		DeparturePortID: req.DeparturePortID,
		ArrivalPortID:   req.ArrivalPortID,
		DepartureDate:   departureDate.Format("2006-01-02"),
		PassengerCount:  req.PassengerCount,
		Limit:           req.Limit,
		Offset:          req.Offset,
	}

	schedules, total, err := s.scheduleRepo.SearchSchedules(ctx, searchReq)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search schedules: %w", err)
	}

	return schedules, total, nil
}

func (s *scheduleService) GetOperatorSchedules(ctx context.Context, operatorID uuid.UUID, date time.Time) ([]*models.Schedule, error) {
	schedules, err := s.scheduleRepo.GetByOperatorAndDate(ctx, operatorID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get operator schedules: %w", err)
	}

	return schedules, nil
}

func (s *scheduleService) GetUpcomingSchedules(ctx context.Context, limit int) ([]*models.Schedule, error) {
	schedules, err := s.scheduleRepo.GetUpcomingSchedules(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming schedules: %w", err)
	}

	return schedules, nil
}