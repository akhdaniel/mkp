package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/ferryflow/boarding-mgt-system/internal/repository"
	"github.com/google/uuid"
)

type BookingService interface {
	CreateBooking(ctx context.Context, customerID uuid.UUID, req *models.CreateBookingRequest) (*models.Booking, error)
	GetBooking(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	GetBookingByReference(ctx context.Context, reference string) (*models.Booking, error)
	CancelBooking(ctx context.Context, id uuid.UUID, reason string) error
	ListBookings(ctx context.Context, filter *models.BookingFilter) ([]*models.Booking, int, error)
	GetCustomerBookings(ctx context.Context, customerID uuid.UUID, limit int) ([]*models.Booking, error)
	GetScheduleManifest(ctx context.Context, scheduleID uuid.UUID) (*models.Manifest, error)
	CheckInTicket(ctx context.Context, qrCode string) error
	GetDailyReport(ctx context.Context, operatorID uuid.UUID, date string) (*models.BookingReport, error)
}

type bookingService struct {
	bookingRepo  repository.BookingRepository
	scheduleRepo repository.ScheduleRepository
	ticketRepo   repository.TicketRepository
	paymentRepo  repository.PaymentRepository
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	scheduleRepo repository.ScheduleRepository,
	ticketRepo repository.TicketRepository,
	paymentRepo repository.PaymentRepository,
) BookingService {
	return &bookingService{
		bookingRepo:  bookingRepo,
		scheduleRepo: scheduleRepo,
		ticketRepo:   ticketRepo,
		paymentRepo:  paymentRepo,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, customerID uuid.UUID, req *models.CreateBookingRequest) (*models.Booking, error) {
	// Get schedule
	schedule, err := s.scheduleRepo.GetByID(ctx, req.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("schedule not found: %w", err)
	}

	// Check if schedule is available
	if schedule.Status != "scheduled" {
		return nil, fmt.Errorf("schedule is not available for booking")
	}

	// Check available seats
	passengerCount := len(req.Passengers)
	if schedule.AvailableSeats < passengerCount {
		return nil, fmt.Errorf("not enough available seats: %d requested, %d available", 
			passengerCount, schedule.AvailableSeats)
	}

	// Calculate total amount
	totalAmount := float64(passengerCount) * schedule.BasePrice

	// Generate booking reference
	bookingRef := s.generateBookingReference()

	// Create booking
	booking := &models.Booking{
		BookingReference:    bookingRef,
		ScheduleID:          req.ScheduleID,
		CustomerID:          customerID,
		PassengerCount:      passengerCount,
		TotalAmount:         totalAmount,
		BookingStatus:       "pending",
		PaymentStatus:       "pending",
		BookingChannel:      "online",
	}

	if req.SpecialRequirements != "" {
		booking.SpecialRequirements = &req.SpecialRequirements
	}

	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// Create tickets for each passenger
	tickets := make([]*models.Ticket, 0, passengerCount)
	for _, passenger := range req.Passengers {
		// Calculate ticket price based on passenger type
		ticketPrice := schedule.BasePrice
		switch passenger.Type {
		case "child":
			ticketPrice = schedule.BasePrice * 0.5
		case "infant":
			ticketPrice = 0
		case "senior":
			ticketPrice = schedule.BasePrice * 0.8
		}

		ticket := &models.Ticket{
			BookingID:      booking.ID,
			PassengerName:  passenger.Name,
			PassengerType:  passenger.Type,
			TicketPrice:    ticketPrice,
			QRCode:         s.generateQRCode(booking.ID, passenger.Name),
			CheckInStatus:  "pending",
		}

		if passenger.SeatNumber != "" {
			ticket.SeatNumber = &passenger.SeatNumber
		}

		tickets = append(tickets, ticket)
	}

	if err := s.ticketRepo.CreateBatch(ctx, tickets); err != nil {
		return nil, fmt.Errorf("failed to create tickets: %w", err)
	}

	// Create payment record
	payment := &models.Payment{
		BookingID:     booking.ID,
		PaymentMethod: req.PaymentMethod,
		Amount:        totalAmount,
		Currency:      "USD",
		PaymentStatus: "pending",
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// TODO: Process payment through gateway

	// For now, simulate successful payment
	booking.BookingStatus = "confirmed"
	booking.PaymentStatus = "completed"
	if err := s.bookingRepo.UpdateStatus(ctx, booking.ID, "confirmed", "completed"); err != nil {
		return nil, fmt.Errorf("failed to update booking status: %w", err)
	}

	// Update payment status
	if err := s.paymentRepo.UpdateStatus(ctx, payment.ID, "completed", nil); err != nil {
		// Non-critical error, log but don't fail
		fmt.Printf("failed to update payment status: %v\n", err)
	}

	// Attach related data
	booking.Schedule = schedule
	booking.Tickets = make([]models.Ticket, len(tickets))
	for i, ticket := range tickets {
		booking.Tickets[i] = *ticket
	}
	booking.Payment = payment

	return booking, nil
}

func (s *bookingService) GetBooking(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("booking not found: %w", err)
	}

	// Get tickets
	tickets, err := s.ticketRepo.GetByBooking(ctx, id)
	if err != nil {
		// Non-critical, continue without tickets
		fmt.Printf("failed to get tickets: %v\n", err)
	} else {
		booking.Tickets = make([]models.Ticket, len(tickets))
		for i, ticket := range tickets {
			booking.Tickets[i] = *ticket
		}
	}

	// Get payment
	payment, err := s.paymentRepo.GetByBooking(ctx, id)
	if err != nil {
		// Non-critical, continue without payment
		fmt.Printf("failed to get payment: %v\n", err)
	} else {
		booking.Payment = payment
	}

	return booking, nil
}

func (s *bookingService) GetBookingByReference(ctx context.Context, reference string) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByReference(ctx, reference)
	if err != nil {
		return nil, fmt.Errorf("booking not found: %w", err)
	}

	return booking, nil
}

func (s *bookingService) CancelBooking(ctx context.Context, id uuid.UUID, reason string) error {
	// Get booking
	booking, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	// Check if booking can be cancelled
	if booking.BookingStatus == "cancelled" {
		return fmt.Errorf("booking is already cancelled")
	}

	// Update booking status
	if err := s.bookingRepo.UpdateStatus(ctx, id, "cancelled", "refund_pending"); err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	// Process refund if payment was completed
	if booking.PaymentStatus == "completed" {
		payment, err := s.paymentRepo.GetByBooking(ctx, id)
		if err == nil && payment != nil {
			if err := s.paymentRepo.ProcessRefund(ctx, payment.ID, payment.Amount); err != nil {
				// Non-critical error, log but don't fail
				fmt.Printf("failed to process refund: %v\n", err)
			}
		}
	}

	return nil
}

func (s *bookingService) ListBookings(ctx context.Context, filter *models.BookingFilter) ([]*models.Booking, int, error) {
	bookings, total, err := s.bookingRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list bookings: %w", err)
	}

	return bookings, total, nil
}

func (s *bookingService) GetCustomerBookings(ctx context.Context, customerID uuid.UUID, limit int) ([]*models.Booking, error) {
	bookings, err := s.bookingRepo.GetCustomerBookings(ctx, customerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer bookings: %w", err)
	}

	return bookings, nil
}

func (s *bookingService) GetScheduleManifest(ctx context.Context, scheduleID uuid.UUID) (*models.Manifest, error) {
	manifest, err := s.ticketRepo.GetManifest(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest: %w", err)
	}

	return manifest, nil
}

func (s *bookingService) CheckInTicket(ctx context.Context, qrCode string) error {
	// Get ticket by QR code
	ticket, err := s.ticketRepo.GetByQRCode(ctx, qrCode)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Check if already checked in
	if ticket.CheckInStatus == "checked_in" {
		return fmt.Errorf("ticket already checked in")
	}

	// Check in the ticket
	if err := s.ticketRepo.CheckIn(ctx, ticket.ID); err != nil {
		return fmt.Errorf("failed to check in ticket: %w", err)
	}

	return nil
}

func (s *bookingService) GetDailyReport(ctx context.Context, operatorID uuid.UUID, date string) (*models.BookingReport, error) {
	report, err := s.bookingRepo.GetDailyReport(ctx, operatorID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily report: %w", err)
	}

	// Parse date for period
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	report.PeriodStart = parsedDate
	report.PeriodEnd = parsedDate.Add(24 * time.Hour)

	return report, nil
}

// Helper functions
func (s *bookingService) generateBookingReference() string {
	// Generate 8-character reference
	b := make([]byte, 6)
	rand.Read(b)
	ref := base64.URLEncoding.EncodeToString(b)
	ref = strings.ToUpper(ref)
	ref = strings.ReplaceAll(ref, "-", "")
	ref = strings.ReplaceAll(ref, "_", "")
	if len(ref) > 8 {
		ref = ref[:8]
	}
	return fmt.Sprintf("FF%s", ref) // FF prefix for FerryFlow
}

func (s *bookingService) generateQRCode(bookingID uuid.UUID, passengerName string) string {
	// Generate unique QR code for ticket
	data := fmt.Sprintf("%s:%s:%d", bookingID.String(), passengerName, time.Now().Unix())
	hash := base64.URLEncoding.EncodeToString([]byte(data))
	return strings.ReplaceAll(hash, "=", "")
}