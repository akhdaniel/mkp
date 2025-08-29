package models

import (
	"time"

	"github.com/google/uuid"
)

// Schedule represents a ferry schedule
type Schedule struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	OperatorID        uuid.UUID  `json:"operator_id" db:"operator_id"`
	RouteID           uuid.UUID  `json:"route_id" db:"route_id"`
	VesselID          uuid.UUID  `json:"vessel_id" db:"vessel_id"`
	DepartureDate     time.Time  `json:"departure_date" db:"departure_date"`
	DepartureTime     time.Time  `json:"departure_time" db:"departure_time"`
	ArrivalTime       time.Time  `json:"arrival_time" db:"arrival_time"`
	BasePrice         float64    `json:"base_price" db:"base_price"`
	TotalCapacity     int        `json:"total_capacity" db:"total_capacity"`
	AvailableSeats    int        `json:"available_seats" db:"available_seats"`
	Status            string     `json:"status" db:"status"`
	CancellationReason *string   `json:"cancellation_reason,omitempty" db:"cancellation_reason"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	Version           int        `json:"version" db:"version"`
	
	// Joined fields
	Operator *Operator `json:"operator,omitempty" db:"-"`
	Route    *Route    `json:"route,omitempty" db:"-"`
	Vessel   *Vessel   `json:"vessel,omitempty" db:"-"`
}

// Booking represents a customer booking
type Booking struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	BookingReference  string     `json:"booking_reference" db:"booking_reference"`
	ScheduleID        uuid.UUID  `json:"schedule_id" db:"schedule_id"`
	CustomerID        uuid.UUID  `json:"customer_id" db:"customer_id"`
	PassengerCount    int        `json:"passenger_count" db:"passenger_count"`
	TotalAmount       float64    `json:"total_amount" db:"total_amount"`
	BookingStatus     string     `json:"booking_status" db:"booking_status"`
	PaymentStatus     string     `json:"payment_status" db:"payment_status"`
	BookingChannel    string     `json:"booking_channel" db:"booking_channel"`
	SpecialRequirements *string  `json:"special_requirements,omitempty" db:"special_requirements"`
	BookingAgentID    *uuid.UUID `json:"booking_agent_id,omitempty" db:"booking_agent_id"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	
	// Joined fields
	Schedule *Schedule `json:"schedule,omitempty" db:"-"`
	Customer *User     `json:"customer,omitempty" db:"-"`
	Tickets  []Ticket  `json:"tickets,omitempty" db:"-"`
	Payment  *Payment  `json:"payment,omitempty" db:"-"`
}

// Ticket represents an individual passenger ticket
type Ticket struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	BookingID     uuid.UUID  `json:"booking_id" db:"booking_id"`
	PassengerName string     `json:"passenger_name" db:"passenger_name"`
	PassengerType string     `json:"passenger_type" db:"passenger_type"`
	SeatNumber    *string    `json:"seat_number,omitempty" db:"seat_number"`
	TicketPrice   float64    `json:"ticket_price" db:"ticket_price"`
	QRCode        string     `json:"qr_code" db:"qr_code"`
	CheckInStatus string     `json:"check_in_status" db:"check_in_status"`
	CheckInTime   *time.Time `json:"check_in_time,omitempty" db:"check_in_time"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	
	// Joined fields
	Booking *Booking `json:"booking,omitempty" db:"-"`
}

// Payment represents a payment transaction
type Payment struct {
	ID                   uuid.UUID              `json:"id" db:"id"`
	BookingID            uuid.UUID              `json:"booking_id" db:"booking_id"`
	PaymentMethod        string                 `json:"payment_method" db:"payment_method"`
	Amount               float64                `json:"amount" db:"amount"`
	Currency             string                 `json:"currency" db:"currency"`
	PaymentStatus        string                 `json:"payment_status" db:"payment_status"`
	GatewayTransactionID *string                `json:"gateway_transaction_id,omitempty" db:"gateway_transaction_id"`
	GatewayResponse      map[string]interface{} `json:"gateway_response,omitempty" db:"gateway_response"`
	ProcessedAt          *time.Time             `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt            time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at" db:"updated_at"`
}

// CreateScheduleRequest represents schedule creation data
type CreateScheduleRequest struct {
	OperatorID    uuid.UUID `json:"operator_id" binding:"required"`
	RouteID       uuid.UUID `json:"route_id" binding:"required"`
	VesselID      uuid.UUID `json:"vessel_id" binding:"required"`
	DepartureDate string    `json:"departure_date" binding:"required"` // Format: "2006-01-02"
	DepartureTime string    `json:"departure_time" binding:"required"` // Format: "15:04"
	ArrivalTime   string    `json:"arrival_time" binding:"required"`   // Format: "15:04"
	BasePrice     float64   `json:"base_price" binding:"required,min=0"`
}

// UpdateScheduleRequest represents schedule update data
type UpdateScheduleRequest struct {
	DepartureDate *string  `json:"departure_date,omitempty"`
	DepartureTime *string  `json:"departure_time,omitempty"`
	ArrivalTime   *string  `json:"arrival_time,omitempty"`
	BasePrice     *float64 `json:"base_price,omitempty"`
	Status        *string  `json:"status,omitempty"`
}

// SearchScheduleRequest represents schedule search criteria
type SearchScheduleRequest struct {
	DeparturePortID uuid.UUID `json:"departure_port_id"`
	ArrivalPortID   uuid.UUID `json:"arrival_port_id"`
	DepartureDate   string    `json:"departure_date"` // Format: "2006-01-02"
	PassengerCount  int       `json:"passenger_count,omitempty"`
	Limit           int       `json:"limit,omitempty"`
	Offset          int       `json:"offset,omitempty"`
}

// CreateBookingRequest represents booking creation data
type CreateBookingRequest struct {
	ScheduleID          uuid.UUID            `json:"schedule_id" binding:"required"`
	Passengers          []PassengerInfo      `json:"passengers" binding:"required,min=1"`
	PaymentMethod       string               `json:"payment_method" binding:"required"`
	SpecialRequirements string               `json:"special_requirements,omitempty"`
}

// PassengerInfo represents passenger information for booking
type PassengerInfo struct {
	Name          string  `json:"name" binding:"required"`
	Type          string  `json:"type" binding:"required,oneof=adult child infant senior"`
	SeatNumber    string  `json:"seat_number,omitempty"`
}

// CancelBookingRequest represents booking cancellation
type CancelBookingRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// CheckInRequest represents ticket check-in
type CheckInRequest struct {
	QRCode string `json:"qr_code" binding:"required"`
}

// BookingFilter represents filters for listing bookings
type BookingFilter struct {
	CustomerID     *uuid.UUID `json:"customer_id,omitempty"`
	ScheduleID     *uuid.UUID `json:"schedule_id,omitempty"`
	BookingStatus  string     `json:"booking_status,omitempty"`
	PaymentStatus  string     `json:"payment_status,omitempty"`
	FromDate       *time.Time `json:"from_date,omitempty"`
	ToDate         *time.Time `json:"to_date,omitempty"`
	Limit          int        `json:"limit,omitempty"`
	Offset         int        `json:"offset,omitempty"`
}

// BookingReport represents booking statistics
type BookingReport struct {
	TotalBookings   int     `json:"total_bookings"`
	TotalRevenue    float64 `json:"total_revenue"`
	TotalPassengers int     `json:"total_passengers"`
	PeriodStart     time.Time `json:"period_start"`
	PeriodEnd       time.Time `json:"period_end"`
	ByStatus        map[string]int `json:"by_status"`
	ByChannel       map[string]int `json:"by_channel"`
}

// RevenueReport represents revenue statistics
type RevenueReport struct {
	TotalRevenue     float64            `json:"total_revenue"`
	RefundedAmount   float64            `json:"refunded_amount"`
	NetRevenue       float64            `json:"net_revenue"`
	PeriodStart      time.Time          `json:"period_start"`
	PeriodEnd        time.Time          `json:"period_end"`
	ByOperator       map[string]float64 `json:"by_operator,omitempty"`
	ByRoute          map[string]float64 `json:"by_route,omitempty"`
	ByPaymentMethod  map[string]float64 `json:"by_payment_method"`
}

// Manifest represents passenger manifest for a schedule
type Manifest struct {
	Schedule       *Schedule `json:"schedule"`
	TotalPassengers int      `json:"total_passengers"`
	CheckedIn      int       `json:"checked_in"`
	Passengers     []ManifestEntry `json:"passengers"`
}

// ManifestEntry represents a passenger entry in manifest
type ManifestEntry struct {
	TicketID       uuid.UUID `json:"ticket_id"`
	PassengerName  string    `json:"passenger_name"`
	PassengerType  string    `json:"passenger_type"`
	SeatNumber     string    `json:"seat_number"`
	BookingRef     string    `json:"booking_ref"`
	CheckInStatus  string    `json:"check_in_status"`
	CustomerEmail  string    `json:"customer_email"`
	CustomerPhone  string    `json:"customer_phone"`
}