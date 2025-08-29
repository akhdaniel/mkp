package service

import (
	"github.com/ferryflow/boarding-mgt-system/internal/auth"
	"github.com/ferryflow/boarding-mgt-system/internal/repository"
)

// Services holds all service interfaces
type Services struct {
	Auth     AuthService
	User     UserService
	Operator OperatorService
	Port     PortService
	Vessel   VesselService
	Route    RouteService
	Schedule ScheduleService
	Booking  BookingService
}

// NewServices creates all service instances
func NewServices(repos *repository.Repositories, jwtUtil *auth.JWTUtil) *Services {
	return &Services{
		Auth:     NewAuthService(repos.User, jwtUtil),
		User:     NewUserService(repos.User),
		Operator: NewOperatorService(repos.Operator),
		Port:     NewPortService(repos.Port),
		Vessel:   NewVesselService(repos.Vessel, repos.Operator),
		Route:    NewRouteService(repos.Route, repos.Port),
		Schedule: NewScheduleService(repos.Schedule, repos.Route, repos.Vessel),
		Booking:  NewBookingService(repos.Booking, repos.Schedule, repos.Ticket, repos.Payment),
	}
}