package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	GetByBooking(ctx context.Context, bookingID uuid.UUID) (*models.Payment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, gatewayTransactionID *string) error
	ProcessRefund(ctx context.Context, paymentID uuid.UUID, amount float64) error
	GetRevenueReport(ctx context.Context, operatorID uuid.UUID, startDate, endDate time.Time) (*models.RevenueReport, error)
}

type paymentRepository struct {
	db *database.DB
}

func NewPaymentRepository(db *database.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	query := `
		INSERT INTO payments (
			booking_id, payment_method, amount, currency,
			payment_status, gateway_transaction_id, gateway_response
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.Pool.QueryRow(ctx, query,
		payment.BookingID, payment.PaymentMethod, payment.Amount,
		payment.Currency, payment.PaymentStatus, payment.GatewayTransactionID,
		payment.GatewayResponse,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	
	return nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	query := `
		SELECT 
			id, booking_id, payment_method, amount, currency,
			payment_status, gateway_transaction_id, gateway_response,
			processed_at, created_at, updated_at
		FROM payments
		WHERE id = $1
	`
	
	payment := &models.Payment{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&payment.ID, &payment.BookingID, &payment.PaymentMethod,
		&payment.Amount, &payment.Currency, &payment.PaymentStatus,
		&payment.GatewayTransactionID, &payment.GatewayResponse,
		&payment.ProcessedAt, &payment.CreatedAt, &payment.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	
	return payment, nil
}

func (r *paymentRepository) GetByBooking(ctx context.Context, bookingID uuid.UUID) (*models.Payment, error) {
	query := `
		SELECT 
			id, booking_id, payment_method, amount, currency,
			payment_status, gateway_transaction_id, gateway_response,
			processed_at, created_at, updated_at
		FROM payments
		WHERE booking_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	
	payment := &models.Payment{}
	err := r.db.Pool.QueryRow(ctx, query, bookingID).Scan(
		&payment.ID, &payment.BookingID, &payment.PaymentMethod,
		&payment.Amount, &payment.Currency, &payment.PaymentStatus,
		&payment.GatewayTransactionID, &payment.GatewayResponse,
		&payment.ProcessedAt, &payment.CreatedAt, &payment.UpdatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	
	return payment, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, gatewayTransactionID *string) error {
	query := `
		UPDATE payments SET
			payment_status = $2,
			gateway_transaction_id = COALESCE($3, gateway_transaction_id),
			processed_at = CASE WHEN $2 IN ('completed', 'failed') THEN CURRENT_TIMESTAMP ELSE processed_at END,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err := r.db.Pool.Exec(ctx, query, id, status, gatewayTransactionID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	
	return nil
}

func (r *paymentRepository) ProcessRefund(ctx context.Context, paymentID uuid.UUID, amount float64) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// Create refund record
	refundQuery := `
		INSERT INTO payments (
			booking_id, payment_method, amount, currency,
			payment_status, gateway_response
		)
		SELECT 
			booking_id, payment_method, -$2, currency,
			'refunded', jsonb_build_object('original_payment_id', $1, 'refund_amount', $2)
		FROM payments
		WHERE id = $1
	`
	
	_, err = tx.Exec(ctx, refundQuery, paymentID, amount)
	if err != nil {
		return fmt.Errorf("failed to create refund record: %w", err)
	}
	
	// Update original payment status
	updateQuery := `
		UPDATE payments SET
			payment_status = 'refunded',
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	_, err = tx.Exec(ctx, updateQuery, paymentID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

func (r *paymentRepository) GetRevenueReport(ctx context.Context, operatorID uuid.UUID, startDate, endDate time.Time) (*models.RevenueReport, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN p.amount > 0 THEN p.amount ELSE 0 END), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN p.amount < 0 THEN ABS(p.amount) ELSE 0 END), 0) as refunded_amount
		FROM payments p
		JOIN bookings b ON p.booking_id = b.id
		JOIN schedules s ON b.schedule_id = s.id
		WHERE s.operator_id = $1 
			AND p.created_at >= $2 
			AND p.created_at <= $3
			AND p.payment_status IN ('completed', 'refunded')
	`
	
	report := &models.RevenueReport{
		PeriodStart:     startDate,
		PeriodEnd:       endDate,
		ByPaymentMethod: make(map[string]float64),
	}
	
	err := r.db.Pool.QueryRow(ctx, query, operatorID, startDate, endDate).Scan(
		&report.TotalRevenue,
		&report.RefundedAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue report: %w", err)
	}
	
	report.NetRevenue = report.TotalRevenue - report.RefundedAmount
	
	// Get breakdown by payment method
	methodQuery := `
		SELECT payment_method, SUM(amount)
		FROM payments p
		JOIN bookings b ON p.booking_id = b.id
		JOIN schedules s ON b.schedule_id = s.id
		WHERE s.operator_id = $1 
			AND p.created_at >= $2 
			AND p.created_at <= $3
			AND p.payment_status = 'completed'
			AND p.amount > 0
		GROUP BY payment_method
	`
	
	methodRows, err := r.db.Pool.Query(ctx, methodQuery, operatorID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get method breakdown: %w", err)
	}
	defer methodRows.Close()
	
	for methodRows.Next() {
		var method string
		var amount float64
		if err := methodRows.Scan(&method, &amount); err != nil {
			return nil, fmt.Errorf("failed to scan method: %w", err)
		}
		report.ByPaymentMethod[method] = amount
	}
	
	return report, nil
}