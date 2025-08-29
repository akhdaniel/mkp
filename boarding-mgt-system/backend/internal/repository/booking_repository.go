package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *models.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error)
	GetByReference(ctx context.Context, reference string) (*models.Booking, error)
	Update(ctx context.Context, booking *models.Booking) error
	UpdateStatus(ctx context.Context, id uuid.UUID, bookingStatus, paymentStatus string) error
	List(ctx context.Context, filter *models.BookingFilter) ([]*models.Booking, int, error)
	GetCustomerBookings(ctx context.Context, customerID uuid.UUID, limit int) ([]*models.Booking, error)
	GetScheduleBookings(ctx context.Context, scheduleID uuid.UUID) ([]*models.Booking, error)
	GetDailyReport(ctx context.Context, operatorID uuid.UUID, date string) (*models.BookingReport, error)
}

type bookingRepository struct {
	db *database.DB
}

func NewBookingRepository(db *database.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// Create booking
	query := `
		INSERT INTO bookings (
			booking_reference, schedule_id, customer_id, passenger_count,
			total_amount, booking_status, payment_status, booking_channel,
			special_requirements, booking_agent_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	
	err = tx.QueryRow(ctx, query,
		booking.BookingReference, booking.ScheduleID, booking.CustomerID,
		booking.PassengerCount, booking.TotalAmount, booking.BookingStatus,
		booking.PaymentStatus, booking.BookingChannel, booking.SpecialRequirements,
		booking.BookingAgentID,
	).Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	
	// Update available seats on schedule (trigger will handle this)
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	query := `
		SELECT 
			b.id, b.booking_reference, b.schedule_id, b.customer_id,
			b.passenger_count, b.total_amount, b.booking_status, b.payment_status,
			b.booking_channel, b.special_requirements, b.booking_agent_id,
			b.created_at, b.updated_at,
			s.id, s.departure_date, s.departure_time, s.arrival_time, s.base_price,
			u.id, u.email, u.first_name, u.last_name, u.phone
		FROM bookings b
		LEFT JOIN schedules s ON b.schedule_id = s.id
		LEFT JOIN users u ON b.customer_id = u.id
		WHERE b.id = $1
	`
	
	booking := &models.Booking{}
	schedule := &models.Schedule{}
	customer := &models.User{}
	
	var specialReq sql.NullString
	var agentID sql.NullString
	var phone sql.NullString
	
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&booking.ID, &booking.BookingReference, &booking.ScheduleID, &booking.CustomerID,
		&booking.PassengerCount, &booking.TotalAmount, &booking.BookingStatus,
		&booking.PaymentStatus, &booking.BookingChannel, &specialReq, &agentID,
		&booking.CreatedAt, &booking.UpdatedAt,
		&schedule.ID, &schedule.DepartureDate, &schedule.DepartureTime, &schedule.ArrivalTime, &schedule.BasePrice,
		&customer.ID, &customer.Email, &customer.FirstName, &customer.LastName, &phone,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("booking not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	
	if specialReq.Valid {
		booking.SpecialRequirements = &specialReq.String
	}
	if agentID.Valid {
		id, _ := uuid.Parse(agentID.String)
		booking.BookingAgentID = &id
	}
	if phone.Valid {
		customer.Phone = &phone.String
	}
	
	booking.Schedule = schedule
	booking.Customer = customer
	
	return booking, nil
}

func (r *bookingRepository) GetByReference(ctx context.Context, reference string) (*models.Booking, error) {
	query := `
		SELECT 
			id, booking_reference, schedule_id, customer_id,
			passenger_count, total_amount, booking_status, payment_status,
			booking_channel, special_requirements, booking_agent_id,
			created_at, updated_at
		FROM bookings
		WHERE booking_reference = $1
	`
	
	booking := &models.Booking{}
	err := r.db.Pool.QueryRow(ctx, query, reference).Scan(
		&booking.ID, &booking.BookingReference, &booking.ScheduleID, &booking.CustomerID,
		&booking.PassengerCount, &booking.TotalAmount, &booking.BookingStatus,
		&booking.PaymentStatus, &booking.BookingChannel, &booking.SpecialRequirements,
		&booking.BookingAgentID, &booking.CreatedAt, &booking.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("booking not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	
	return booking, nil
}

func (r *bookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	query := `
		UPDATE bookings SET
			passenger_count = $2,
			total_amount = $3,
			booking_status = $4,
			payment_status = $5,
			special_requirements = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		booking.ID, booking.PassengerCount, booking.TotalAmount,
		booking.BookingStatus, booking.PaymentStatus, booking.SpecialRequirements,
	).Scan(&booking.UpdatedAt)
	
	if err == pgx.ErrNoRows {
		return fmt.Errorf("booking not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	
	return nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, bookingStatus, paymentStatus string) error {
	query := `
		UPDATE bookings SET
			booking_status = $2,
			payment_status = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err := r.db.Pool.Exec(ctx, query, id, bookingStatus, paymentStatus)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}
	
	return nil
}

func (r *bookingRepository) List(ctx context.Context, filter *models.BookingFilter) ([]*models.Booking, int, error) {
	query := `
		SELECT 
			id, booking_reference, schedule_id, customer_id,
			passenger_count, total_amount, booking_status, payment_status,
			booking_channel, special_requirements, booking_agent_id,
			created_at, updated_at
		FROM bookings
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM bookings WHERE 1=1`
	
	args := []interface{}{}
	argCount := 0
	
	// Build filters
	if filter.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		countQuery += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filter.CustomerID)
	}
	
	if filter.ScheduleID != nil {
		argCount++
		query += fmt.Sprintf(" AND schedule_id = $%d", argCount)
		countQuery += fmt.Sprintf(" AND schedule_id = $%d", argCount)
		args = append(args, *filter.ScheduleID)
	}
	
	if filter.BookingStatus != "" {
		argCount++
		query += fmt.Sprintf(" AND booking_status = $%d", argCount)
		countQuery += fmt.Sprintf(" AND booking_status = $%d", argCount)
		args = append(args, filter.BookingStatus)
	}
	
	if filter.PaymentStatus != "" {
		argCount++
		query += fmt.Sprintf(" AND payment_status = $%d", argCount)
		countQuery += fmt.Sprintf(" AND payment_status = $%d", argCount)
		args = append(args, filter.PaymentStatus)
	}
	
	if filter.FromDate != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		countQuery += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, *filter.FromDate)
	}
	
	if filter.ToDate != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		countQuery += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, *filter.ToDate)
	}
	
	// Get total count
	var totalCount int
	err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}
	
	// Add pagination
	query += " ORDER BY created_at DESC"
	if filter.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}
	
	// Execute query
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list bookings: %w", err)
	}
	defer rows.Close()
	
	bookings := []*models.Booking{}
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.BookingReference, &booking.ScheduleID, &booking.CustomerID,
			&booking.PassengerCount, &booking.TotalAmount, &booking.BookingStatus,
			&booking.PaymentStatus, &booking.BookingChannel, &booking.SpecialRequirements,
			&booking.BookingAgentID, &booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}
	
	return bookings, totalCount, nil
}

func (r *bookingRepository) GetCustomerBookings(ctx context.Context, customerID uuid.UUID, limit int) ([]*models.Booking, error) {
	query := `
		SELECT 
			b.id, b.booking_reference, b.schedule_id, b.customer_id,
			b.passenger_count, b.total_amount, b.booking_status, b.payment_status,
			b.booking_channel, b.special_requirements, b.booking_agent_id,
			b.created_at, b.updated_at
		FROM bookings b
		JOIN schedules s ON b.schedule_id = s.id
		WHERE b.customer_id = $1
		ORDER BY s.departure_date DESC, s.departure_time DESC
		LIMIT $2
	`
	
	rows, err := r.db.Pool.Query(ctx, query, customerID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer bookings: %w", err)
	}
	defer rows.Close()
	
	bookings := []*models.Booking{}
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.BookingReference, &booking.ScheduleID, &booking.CustomerID,
			&booking.PassengerCount, &booking.TotalAmount, &booking.BookingStatus,
			&booking.PaymentStatus, &booking.BookingChannel, &booking.SpecialRequirements,
			&booking.BookingAgentID, &booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}
	
	return bookings, nil
}

func (r *bookingRepository) GetScheduleBookings(ctx context.Context, scheduleID uuid.UUID) ([]*models.Booking, error) {
	query := `
		SELECT 
			id, booking_reference, schedule_id, customer_id,
			passenger_count, total_amount, booking_status, payment_status,
			booking_channel, special_requirements, booking_agent_id,
			created_at, updated_at
		FROM bookings
		WHERE schedule_id = $1 AND booking_status = 'confirmed'
		ORDER BY created_at ASC
	`
	
	rows, err := r.db.Pool.Query(ctx, query, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule bookings: %w", err)
	}
	defer rows.Close()
	
	bookings := []*models.Booking{}
	for rows.Next() {
		booking := &models.Booking{}
		err := rows.Scan(
			&booking.ID, &booking.BookingReference, &booking.ScheduleID, &booking.CustomerID,
			&booking.PassengerCount, &booking.TotalAmount, &booking.BookingStatus,
			&booking.PaymentStatus, &booking.BookingChannel, &booking.SpecialRequirements,
			&booking.BookingAgentID, &booking.CreatedAt, &booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}
	
	return bookings, nil
}

func (r *bookingRepository) GetDailyReport(ctx context.Context, operatorID uuid.UUID, date string) (*models.BookingReport, error) {
	query := `
		SELECT 
			COUNT(*) as total_bookings,
			COALESCE(SUM(b.total_amount), 0) as total_revenue,
			COALESCE(SUM(b.passenger_count), 0) as total_passengers
		FROM bookings b
		JOIN schedules s ON b.schedule_id = s.id
		WHERE s.operator_id = $1 
			AND DATE(b.created_at) = $2
			AND b.booking_status = 'confirmed'
	`
	
	report := &models.BookingReport{
		ByStatus:  make(map[string]int),
		ByChannel: make(map[string]int),
	}
	
	err := r.db.Pool.QueryRow(ctx, query, operatorID, date).Scan(
		&report.TotalBookings,
		&report.TotalRevenue,
		&report.TotalPassengers,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking report: %w", err)
	}
	
	// Get breakdown by status
	statusQuery := `
		SELECT booking_status, COUNT(*)
		FROM bookings b
		JOIN schedules s ON b.schedule_id = s.id
		WHERE s.operator_id = $1 AND DATE(b.created_at) = $2
		GROUP BY booking_status
	`
	
	statusRows, err := r.db.Pool.Query(ctx, statusQuery, operatorID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get status breakdown: %w", err)
	}
	defer statusRows.Close()
	
	for statusRows.Next() {
		var status string
		var count int
		if err := statusRows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan status: %w", err)
		}
		report.ByStatus[status] = count
	}
	
	// Get breakdown by channel
	channelQuery := `
		SELECT booking_channel, COUNT(*)
		FROM bookings b
		JOIN schedules s ON b.schedule_id = s.id
		WHERE s.operator_id = $1 AND DATE(b.created_at) = $2
		GROUP BY booking_channel
	`
	
	channelRows, err := r.db.Pool.Query(ctx, channelQuery, operatorID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel breakdown: %w", err)
	}
	defer channelRows.Close()
	
	for channelRows.Next() {
		var channel string
		var count int
		if err := channelRows.Scan(&channel, &count); err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		report.ByChannel[channel] = count
	}
	
	return report, nil
}