package repository

import (
	"context"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TicketRepository interface {
	CreateBatch(ctx context.Context, tickets []*models.Ticket) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error)
	GetByQRCode(ctx context.Context, qrCode string) (*models.Ticket, error)
	GetByBooking(ctx context.Context, bookingID uuid.UUID) ([]*models.Ticket, error)
	CheckIn(ctx context.Context, ticketID uuid.UUID) error
	GetManifest(ctx context.Context, scheduleID uuid.UUID) (*models.Manifest, error)
}

type ticketRepository struct {
	db *database.DB
}

func NewTicketRepository(db *database.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) CreateBatch(ctx context.Context, tickets []*models.Ticket) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	query := `
		INSERT INTO tickets (
			booking_id, passenger_name, passenger_type, seat_number,
			ticket_price, qr_code, check_in_status
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	
	for _, ticket := range tickets {
		err := tx.QueryRow(ctx, query,
			ticket.BookingID, ticket.PassengerName, ticket.PassengerType,
			ticket.SeatNumber, ticket.TicketPrice, ticket.QRCode, ticket.CheckInStatus,
		).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt)
		
		if err != nil {
			return fmt.Errorf("failed to create ticket: %w", err)
		}
	}
	
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

func (r *ticketRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Ticket, error) {
	query := `
		SELECT 
			t.id, t.booking_id, t.passenger_name, t.passenger_type,
			t.seat_number, t.ticket_price, t.qr_code, t.check_in_status,
			t.check_in_time, t.created_at, t.updated_at,
			b.id, b.booking_reference, b.schedule_id, b.customer_id,
			b.total_amount, b.booking_status
		FROM tickets t
		LEFT JOIN bookings b ON t.booking_id = b.id
		WHERE t.id = $1
	`
	
	ticket := &models.Ticket{}
	booking := &models.Booking{}
	
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&ticket.ID, &ticket.BookingID, &ticket.PassengerName, &ticket.PassengerType,
		&ticket.SeatNumber, &ticket.TicketPrice, &ticket.QRCode, &ticket.CheckInStatus,
		&ticket.CheckInTime, &ticket.CreatedAt, &ticket.UpdatedAt,
		&booking.ID, &booking.BookingReference, &booking.ScheduleID, &booking.CustomerID,
		&booking.TotalAmount, &booking.BookingStatus,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("ticket not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	
	ticket.Booking = booking
	return ticket, nil
}

func (r *ticketRepository) GetByQRCode(ctx context.Context, qrCode string) (*models.Ticket, error) {
	query := `
		SELECT 
			id, booking_id, passenger_name, passenger_type,
			seat_number, ticket_price, qr_code, check_in_status,
			check_in_time, created_at, updated_at
		FROM tickets
		WHERE qr_code = $1
	`
	
	ticket := &models.Ticket{}
	err := r.db.Pool.QueryRow(ctx, query, qrCode).Scan(
		&ticket.ID, &ticket.BookingID, &ticket.PassengerName, &ticket.PassengerType,
		&ticket.SeatNumber, &ticket.TicketPrice, &ticket.QRCode, &ticket.CheckInStatus,
		&ticket.CheckInTime, &ticket.CreatedAt, &ticket.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("ticket not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	
	return ticket, nil
}

func (r *ticketRepository) GetByBooking(ctx context.Context, bookingID uuid.UUID) ([]*models.Ticket, error) {
	query := `
		SELECT 
			id, booking_id, passenger_name, passenger_type,
			seat_number, ticket_price, qr_code, check_in_status,
			check_in_time, created_at, updated_at
		FROM tickets
		WHERE booking_id = $1
		ORDER BY created_at ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, query, bookingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	defer rows.Close()
	
	tickets := []*models.Ticket{}
	for rows.Next() {
		ticket := &models.Ticket{}
		err := rows.Scan(
			&ticket.ID, &ticket.BookingID, &ticket.PassengerName, &ticket.PassengerType,
			&ticket.SeatNumber, &ticket.TicketPrice, &ticket.QRCode, &ticket.CheckInStatus,
			&ticket.CheckInTime, &ticket.CreatedAt, &ticket.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}
		tickets = append(tickets, ticket)
	}
	
	return tickets, nil
}

func (r *ticketRepository) CheckIn(ctx context.Context, ticketID uuid.UUID) error {
	query := `
		UPDATE tickets SET
			check_in_status = 'checked_in',
			check_in_time = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND check_in_status = 'pending'
	`
	
	result, err := r.db.Pool.Exec(ctx, query, ticketID)
	if err != nil {
		return fmt.Errorf("failed to check in ticket: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("ticket not found or already checked in")
	}
	
	return nil
}

func (r *ticketRepository) GetManifest(ctx context.Context, scheduleID uuid.UUID) (*models.Manifest, error) {
	// Get schedule details
	scheduleQuery := `
		SELECT 
			s.id, s.departure_date, s.departure_time, s.arrival_time,
			s.total_capacity, s.available_seats
		FROM schedules s
		WHERE s.id = $1
	`
	
	manifest := &models.Manifest{
		Schedule: &models.Schedule{},
	}
	
	err := r.db.Pool.QueryRow(ctx, scheduleQuery, scheduleID).Scan(
		&manifest.Schedule.ID, &manifest.Schedule.DepartureDate,
		&manifest.Schedule.DepartureTime, &manifest.Schedule.ArrivalTime,
		&manifest.Schedule.TotalCapacity, &manifest.Schedule.AvailableSeats,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}
	
	// Get passenger details
	passengerQuery := `
		SELECT 
			t.id, t.passenger_name, t.passenger_type, t.seat_number,
			b.booking_reference, t.check_in_status,
			u.email, u.phone
		FROM tickets t
		JOIN bookings b ON t.booking_id = b.id
		JOIN users u ON b.customer_id = u.id
		WHERE b.schedule_id = $1 AND b.booking_status = 'confirmed'
		ORDER BY t.seat_number ASC, t.passenger_name ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, passengerQuery, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get passengers: %w", err)
	}
	defer rows.Close()
	
	manifest.Passengers = []models.ManifestEntry{}
	checkedInCount := 0
	
	for rows.Next() {
		entry := models.ManifestEntry{}
		var seatNumber, phone *string
		
		err := rows.Scan(
			&entry.TicketID, &entry.PassengerName, &entry.PassengerType,
			&seatNumber, &entry.BookingRef, &entry.CheckInStatus,
			&entry.CustomerEmail, &phone,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan passenger: %w", err)
		}
		
		if seatNumber != nil {
			entry.SeatNumber = *seatNumber
		}
		if phone != nil {
			entry.CustomerPhone = *phone
		}
		
		if entry.CheckInStatus == "checked_in" {
			checkedInCount++
		}
		
		manifest.Passengers = append(manifest.Passengers, entry)
	}
	
	manifest.TotalPassengers = len(manifest.Passengers)
	manifest.CheckedIn = checkedInCount
	
	return manifest, nil
}