package repository

import (
	"github.com/ferryflow/boarding-mgt-system/internal/database"
)

// Repositories holds all repository interfaces
type Repositories struct {
	User     UserRepository
	Operator OperatorRepository
	Port     PortRepository
	Vessel   VesselRepository
	Route    RouteRepository
	Schedule ScheduleRepository
	Booking  BookingRepository
	Ticket   TicketRepository
	Payment  PaymentRepository
}

// NewRepositories creates all repository instances
func NewRepositories(db *database.DB) *Repositories {
	return &Repositories{
		User:     NewUserRepository(db),
		Operator: NewOperatorRepository(db),
		Port:     NewPortRepository(db),
		Vessel:   NewVesselRepository(db),
		Route:    NewRouteRepository(db),
		Schedule: NewScheduleRepository(db),
		Booking:  NewBookingRepository(db),
		Ticket:   NewTicketRepository(db),
		Payment:  NewPaymentRepository(db),
	}
}