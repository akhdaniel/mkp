package models

import (
	"time"

	"github.com/google/uuid"
)

// Operator represents a ferry operator company
type Operator struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	Name         string            `json:"name" db:"name"`
	Code         string            `json:"code" db:"code"`
	ContactEmail string            `json:"contact_email" db:"contact_email"`
	ContactPhone *string           `json:"contact_phone,omitempty" db:"contact_phone"`
	Address      *string           `json:"address,omitempty" db:"address"`
	IsActive     bool              `json:"is_active" db:"is_active"`
	Settings     map[string]interface{} `json:"settings" db:"settings"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
}

// Port represents a ferry terminal/port
type Port struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Code        string                 `json:"code" db:"code"`
	City        string                 `json:"city" db:"city"`
	Country     string                 `json:"country" db:"country"`
	Timezone    string                 `json:"timezone" db:"timezone"`
	Coordinates *Coordinates           `json:"coordinates,omitempty" db:"coordinates"`
	Facilities  map[string]interface{} `json:"facilities" db:"facilities"`
	IsActive    bool                   `json:"is_active" db:"is_active"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// Coordinates represents geographic coordinates
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Vessel represents a ferry vessel
type Vessel struct {
	ID               uuid.UUID              `json:"id" db:"id"`
	OperatorID       uuid.UUID              `json:"operator_id" db:"operator_id"`
	Name             string                 `json:"name" db:"name"`
	RegistrationNumber string               `json:"registration_number" db:"registration_number"`
	VesselType       string                 `json:"vessel_type" db:"vessel_type"`
	Capacity         int                    `json:"capacity" db:"capacity"`
	DeckCount        int                    `json:"deck_count" db:"deck_count"`
	SeatConfiguration map[string]interface{} `json:"seat_configuration" db:"seat_configuration"`
	Amenities        map[string]interface{} `json:"amenities" db:"amenities"`
	IsActive         bool                   `json:"is_active" db:"is_active"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
	
	// Joined fields
	Operator *Operator `json:"operator,omitempty" db:"-"`
}

// Route represents a route between two ports
type Route struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	OperatorID      uuid.UUID     `json:"operator_id" db:"operator_id"`
	Name            string        `json:"name" db:"name"`
	DeparturePortID uuid.UUID     `json:"departure_port_id" db:"departure_port_id"`
	ArrivalPortID   uuid.UUID     `json:"arrival_port_id" db:"arrival_port_id"`
	DistanceKM      *float64      `json:"distance_km,omitempty" db:"distance_km"`
	EstimatedDuration time.Duration `json:"estimated_duration" db:"estimated_duration"`
	IsActive        bool          `json:"is_active" db:"is_active"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	
	// Joined fields
	Operator      *Operator `json:"operator,omitempty" db:"-"`
	DeparturePort *Port     `json:"departure_port,omitempty" db:"-"`
	ArrivalPort   *Port     `json:"arrival_port,omitempty" db:"-"`
}

// CreateOperatorRequest represents operator creation data
type CreateOperatorRequest struct {
	Name         string                 `json:"name" binding:"required"`
	Code         string                 `json:"code" binding:"required,min=3,max=10"`
	ContactEmail string                 `json:"contact_email" binding:"required,email"`
	ContactPhone string                 `json:"contact_phone,omitempty"`
	Address      string                 `json:"address,omitempty"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
}

// UpdateOperatorRequest represents operator update data
type UpdateOperatorRequest struct {
	Name         *string                `json:"name,omitempty"`
	ContactEmail *string                `json:"contact_email,omitempty"`
	ContactPhone *string                `json:"contact_phone,omitempty"`
	Address      *string                `json:"address,omitempty"`
	IsActive     *bool                  `json:"is_active,omitempty"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
}

// CreatePortRequest represents port creation data
type CreatePortRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Code        string                 `json:"code" binding:"required,min=3,max=10"`
	City        string                 `json:"city" binding:"required"`
	Country     string                 `json:"country" binding:"required"`
	Timezone    string                 `json:"timezone" binding:"required"`
	Coordinates *Coordinates           `json:"coordinates,omitempty"`
	Facilities  map[string]interface{} `json:"facilities,omitempty"`
}

// UpdatePortRequest represents port update data
type UpdatePortRequest struct {
	Name        *string                `json:"name,omitempty"`
	City        *string                `json:"city,omitempty"`
	Country     *string                `json:"country,omitempty"`
	Timezone    *string                `json:"timezone,omitempty"`
	Coordinates *Coordinates           `json:"coordinates,omitempty"`
	Facilities  map[string]interface{} `json:"facilities,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
}

// CreateVesselRequest represents vessel creation data
type CreateVesselRequest struct {
	OperatorID         uuid.UUID              `json:"operator_id" binding:"required"`
	Name               string                 `json:"name" binding:"required"`
	RegistrationNumber string                 `json:"registration_number" binding:"required"`
	VesselType         string                 `json:"vessel_type" binding:"required,oneof=passenger cargo mixed"`
	Capacity           int                    `json:"capacity" binding:"required,min=1"`
	DeckCount          int                    `json:"deck_count" binding:"required,min=1"`
	SeatConfiguration  map[string]interface{} `json:"seat_configuration" binding:"required"`
	Amenities          map[string]interface{} `json:"amenities,omitempty"`
}

// UpdateVesselRequest represents vessel update data
type UpdateVesselRequest struct {
	Name              *string                `json:"name,omitempty"`
	VesselType        *string                `json:"vessel_type,omitempty"`
	Capacity          *int                   `json:"capacity,omitempty"`
	DeckCount         *int                   `json:"deck_count,omitempty"`
	SeatConfiguration map[string]interface{} `json:"seat_configuration,omitempty"`
	Amenities         map[string]interface{} `json:"amenities,omitempty"`
	IsActive          *bool                  `json:"is_active,omitempty"`
}

// CreateRouteRequest represents route creation data
type CreateRouteRequest struct {
	OperatorID        uuid.UUID `json:"operator_id" binding:"required"`
	Name              string    `json:"name" binding:"required"`
	DeparturePortID   uuid.UUID `json:"departure_port_id" binding:"required"`
	ArrivalPortID     uuid.UUID `json:"arrival_port_id" binding:"required"`
	DistanceKM        float64   `json:"distance_km,omitempty"`
	EstimatedDuration string    `json:"estimated_duration" binding:"required"` // Format: "2h30m"
}

// UpdateRouteRequest represents route update data
type UpdateRouteRequest struct {
	Name              *string  `json:"name,omitempty"`
	DistanceKM        *float64 `json:"distance_km,omitempty"`
	EstimatedDuration *string  `json:"estimated_duration,omitempty"`
	IsActive          *bool    `json:"is_active,omitempty"`
}