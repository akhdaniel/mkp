package seed

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Seeder struct {
	db *pgxpool.Pool
}

func NewSeeder(db *pgxpool.Pool) *Seeder {
	return &Seeder{db: db}
}

// SeedAll seeds all data for development
func (s *Seeder) SeedAll(ctx context.Context) error {
	// Seed in order of dependencies
	if err := s.SeedOperators(ctx); err != nil {
		return fmt.Errorf("failed to seed operators: %w", err)
	}

	if err := s.SeedPorts(ctx); err != nil {
		return fmt.Errorf("failed to seed ports: %w", err)
	}

	if err := s.SeedVessels(ctx); err != nil {
		return fmt.Errorf("failed to seed vessels: %w", err)
	}

	if err := s.SeedRoutes(ctx); err != nil {
		return fmt.Errorf("failed to seed routes: %w", err)
	}

	if err := s.SeedUsers(ctx); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	if err := s.SeedSchedules(ctx); err != nil {
		return fmt.Errorf("failed to seed schedules: %w", err)
	}

	if err := s.SeedBookings(ctx); err != nil {
		return fmt.Errorf("failed to seed bookings: %w", err)
	}

	return nil
}

// SeedOperators creates sample ferry operators
func (s *Seeder) SeedOperators(ctx context.Context) error {
	operators := []struct {
		name    string
		code    string
		email   string
		phone   string
		address string
	}{
		{"Island Express Ferries", "IEF001", "contact@islandexpress.com", "+1234567890", "123 Harbor Way, Port City"},
		{"Coastal Connect Lines", "CCL001", "info@coastalconnect.com", "+1234567891", "456 Marina Blvd, Seaside"},
		{"Blue Wave Transport", "BWT001", "support@bluewave.com", "+1234567892", "789 Ocean Drive, Bay Town"},
	}

	for _, op := range operators {
		_, err := s.db.Exec(ctx, `
			INSERT INTO operators (name, code, contact_email, contact_phone, address)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (code) DO NOTHING
		`, op.name, op.code, op.email, op.phone, op.address)
		if err != nil {
			return err
		}
	}

	return nil
}

// SeedPorts creates sample ports
func (s *Seeder) SeedPorts(ctx context.Context) error {
	ports := []struct {
		name     string
		code     string
		city     string
		country  string
		timezone string
		lat      float64
		lon      float64
	}{
		{"Main Harbor Terminal", "MHT", "Harbor City", "USA", "America/New_York", 40.7128, -74.0060},
		{"Northport Ferry Terminal", "NPT", "Northport", "USA", "America/New_York", 40.9006, -73.3439},
		{"Southbay Marina", "SBM", "South Bay", "USA", "America/New_York", 40.5897, -74.1915},
		{"Island Cove Port", "ICP", "Island Cove", "USA", "America/New_York", 41.0534, -73.5387},
		{"Eastside Dock", "ESD", "East Harbor", "USA", "America/New_York", 40.7614, -73.9776},
	}

	for _, p := range ports {
		_, err := s.db.Exec(ctx, `
			INSERT INTO ports (name, code, city, country, timezone, coordinates)
			VALUES ($1, $2, $3, $4, $5, POINT($6, $7))
			ON CONFLICT (code) DO NOTHING
		`, p.name, p.code, p.city, p.country, p.timezone, p.lat, p.lon)
		if err != nil {
			return err
		}
	}

	return nil
}

// SeedVessels creates sample vessels
func (s *Seeder) SeedVessels(ctx context.Context) error {
	// Get operator IDs
	rows, err := s.db.Query(ctx, "SELECT id, code FROM operators")
	if err != nil {
		return err
	}
	defer rows.Close()

	operatorMap := make(map[string]string)
	for rows.Next() {
		var id, code string
		if err := rows.Scan(&id, &code); err != nil {
			return err
		}
		operatorMap[code] = id
	}

	vessels := []struct {
		operatorCode string
		name         string
		regNumber    string
		vesselType   string
		capacity     int
		deckCount    int
	}{
		{"IEF001", "MV Island Explorer", "VES001", "passenger", 250, 2},
		{"IEF001", "MV Sea Breeze", "VES002", "passenger", 180, 2},
		{"CCL001", "MV Coastal Runner", "VES003", "passenger", 300, 3},
		{"CCL001", "MV Harbor Queen", "VES004", "mixed", 200, 2},
		{"BWT001", "MV Blue Horizon", "VES005", "passenger", 150, 1},
		{"BWT001", "MV Wave Rider", "VES006", "passenger", 220, 2},
	}

	for _, v := range vessels {
		operatorID := operatorMap[v.operatorCode]
		seatConfig := fmt.Sprintf(`{"decks": %d, "seats_per_deck": %d}`, v.deckCount, v.capacity/v.deckCount)
		
		_, err := s.db.Exec(ctx, `
			INSERT INTO vessels (operator_id, name, registration_number, vessel_type, capacity, deck_count, seat_configuration)
			VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb)
			ON CONFLICT (registration_number) DO NOTHING
		`, operatorID, v.name, v.regNumber, v.vesselType, v.capacity, v.deckCount, seatConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

// SeedRoutes creates sample routes
func (s *Seeder) SeedRoutes(ctx context.Context) error {
	// Get operator and port IDs
	var operators []struct{ id, code string }
	var ports []struct{ id, code string }

	rows, _ := s.db.Query(ctx, "SELECT id, code FROM operators")
	for rows.Next() {
		var op struct{ id, code string }
		rows.Scan(&op.id, &op.code)
		operators = append(operators, op)
	}
	rows.Close()

	rows, _ = s.db.Query(ctx, "SELECT id, code FROM ports")
	for rows.Next() {
		var p struct{ id, code string }
		rows.Scan(&p.id, &p.code)
		ports = append(ports, p)
	}
	rows.Close()

	// Create routes between ports
	if len(operators) > 0 && len(ports) >= 2 {
		routes := []struct {
			name     string
			depPort  int
			arrPort  int
			distance float64
			duration string
		}{
			{"Harbor to Northport Express", 0, 1, 45.5, "1 hour 30 mins"},
			{"Harbor to South Bay", 0, 2, 32.0, "1 hour"},
			{"Northport to Island Cove", 1, 3, 28.5, "45 mins"},
			{"South Bay to Eastside", 2, 4, 38.0, "1 hour 15 mins"},
			{"Island Cove to Harbor", 3, 0, 52.0, "1 hour 45 mins"},
		}

		for i, r := range routes {
			operatorID := operators[i%len(operators)].id
			depPortID := ports[r.depPort].id
			arrPortID := ports[r.arrPort].id

			_, err := s.db.Exec(ctx, `
				INSERT INTO routes (operator_id, name, departure_port_id, arrival_port_id, distance_km, estimated_duration)
				VALUES ($1, $2, $3, $4, $5, $6::interval)
				ON CONFLICT DO NOTHING
			`, operatorID, r.name, depPortID, arrPortID, r.distance, r.duration)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedUsers creates sample users
func (s *Seeder) SeedUsers(ctx context.Context) error {
	// Get first operator ID for agents
	var operatorID string
	err := s.db.QueryRow(ctx, "SELECT id FROM operators LIMIT 1").Scan(&operatorID)
	if err != nil {
		return err
	}

	// Default password for all test users: "Password123!"
	passwordHash, _ := auth.HashPassword("Password123!")

	users := []struct {
		email      string
		firstName  string
		lastName   string
		phone      string
		userType   string
		operatorID *string
	}{
		// Customers
		{"john.doe@example.com", "John", "Doe", "+1234567800", "customer", nil},
		{"jane.smith@example.com", "Jane", "Smith", "+1234567801", "customer", nil},
		{"robert.jones@example.com", "Robert", "Jones", "+1234567802", "customer", nil},
		{"maria.garcia@example.com", "Maria", "Garcia", "+1234567803", "customer", nil},
		// Agents
		{"agent.wilson@ferryops.com", "Alex", "Wilson", "+1234567810", "agent", &operatorID},
		{"agent.brown@ferryops.com", "Sarah", "Brown", "+1234567811", "agent", &operatorID},
		// Operator Admin
		{"admin@ferryops.com", "Admin", "User", "+1234567820", "operator_admin", &operatorID},
		// System Admin
		{"system@ferryflow.com", "System", "Admin", "+1234567830", "system_admin", nil},
	}

	for _, u := range users {
		_, err := s.db.Exec(ctx, `
			INSERT INTO users (email, password_hash, first_name, last_name, phone, user_type, operator_id, is_verified, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, true, true)
			ON CONFLICT (email) DO NOTHING
		`, u.email, passwordHash, u.firstName, u.lastName, u.phone, u.userType, u.operatorID)
		if err != nil {
			return err
		}
	}

	return nil
}

// SeedSchedules creates sample schedules
func (s *Seeder) SeedSchedules(ctx context.Context) error {
	// Get routes and vessels
	rows, err := s.db.Query(ctx, `
		SELECT r.id, r.operator_id, v.id as vessel_id
		FROM routes r
		JOIN vessels v ON v.operator_id = r.operator_id
		LIMIT 10
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	type routeVessel struct {
		routeID    string
		operatorID string
		vesselID   string
	}
	var routeVessels []routeVessel

	for rows.Next() {
		var rv routeVessel
		if err := rows.Scan(&rv.routeID, &rv.operatorID, &rv.vesselID); err != nil {
			return err
		}
		routeVessels = append(routeVessels, rv)
	}

	// Create schedules for the next 30 days
	now := time.Now()
	for _, rv := range routeVessels {
		for days := 0; days < 30; days++ {
			date := now.AddDate(0, 0, days)
			
			// Morning and afternoon departures
			departures := []string{"08:00:00", "14:00:00", "18:00:00"}
			arrivals := []string{"10:00:00", "16:00:00", "20:00:00"}
			prices := []float64{35.00, 40.00, 35.00}

			for i, depTime := range departures {
				// Random capacity between 100-200
				capacity := 100 + rand.Intn(100)
				available := capacity - rand.Intn(20) // Some bookings already made

				_, err := s.db.Exec(ctx, `
					INSERT INTO schedules (
						operator_id, route_id, vessel_id,
						departure_date, departure_time, arrival_time,
						base_price, total_capacity, available_seats, status
					) VALUES ($1, $2, $3, $4, $5::time, $6::time, $7, $8, $9, 'scheduled')
					ON CONFLICT DO NOTHING
				`, rv.operatorID, rv.routeID, rv.vesselID,
					date.Format("2006-01-02"), depTime, arrivals[i],
					prices[i], capacity, available)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// SeedBookings creates sample bookings
func (s *Seeder) SeedBookings(ctx context.Context) error {
	// Get customers and schedules
	var customers []string
	rows, _ := s.db.Query(ctx, "SELECT id FROM users WHERE user_type = 'customer' LIMIT 10")
	for rows.Next() {
		var id string
		rows.Scan(&id)
		customers = append(customers, id)
	}
	rows.Close()

	var schedules []string
	rows, _ = s.db.Query(ctx, `
		SELECT id FROM schedules 
		WHERE departure_date >= CURRENT_DATE 
		AND available_seats > 10 
		LIMIT 20
	`)
	for rows.Next() {
		var id string
		rows.Scan(&id)
		schedules = append(schedules, id)
	}
	rows.Close()

	// Create bookings
	for i := 0; i < len(schedules) && i < len(customers)*2; i++ {
		customerID := customers[i%len(customers)]
		scheduleID := schedules[i]
		passengerCount := 1 + rand.Intn(4)
		totalAmount := float64(passengerCount) * (35.0 + float64(rand.Intn(15)))

		bookingRef := fmt.Sprintf("BK%d%04d", time.Now().Unix(), i)

		bookingID := ""
		err := s.db.QueryRow(ctx, `
			INSERT INTO bookings (
				booking_reference, schedule_id, customer_id,
				passenger_count, total_amount, booking_status,
				payment_status, booking_channel
			) VALUES ($1, $2, $3, $4, $5, 'confirmed', 'paid', 'online')
			ON CONFLICT (booking_reference) DO NOTHING
			RETURNING id
		`, bookingRef, scheduleID, customerID, passengerCount, totalAmount).Scan(&bookingID)

		if err == nil && bookingID != "" {
			// Create tickets for the booking
			for j := 0; j < passengerCount; j++ {
				qrCode := fmt.Sprintf("QR-%s-%d", bookingRef, j)
				passengerName := fmt.Sprintf("Passenger %d", j+1)
				seatNumber := fmt.Sprintf("%c%d", 'A'+j/10, (j%10)+1)

				_, err := s.db.Exec(ctx, `
					INSERT INTO tickets (
						booking_id, passenger_name, passenger_type,
						seat_number, ticket_price, qr_code
					) VALUES ($1, $2, 'adult', $3, $4, $5)
					ON CONFLICT (qr_code) DO NOTHING
				`, bookingID, passengerName, seatNumber, totalAmount/float64(passengerCount), qrCode)
				if err != nil {
					return err
				}
			}

			// Create payment record
			_, err = s.db.Exec(ctx, `
				INSERT INTO payments (
					booking_id, payment_method, amount,
					payment_status, processed_at
				) VALUES ($1, 'credit_card', $2, 'completed', CURRENT_TIMESTAMP)
			`, bookingID, totalAmount)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CleanAll removes all seeded data
func (s *Seeder) CleanAll(ctx context.Context) error {
	// Delete in reverse order of dependencies
	tables := []string{
		"support_messages",
		"support_tickets",
		"refunds",
		"payments",
		"tickets",
		"bookings",
		"schedules",
		"user_sessions",
		"users",
		"routes",
		"vessels",
		"ports",
		"operators",
	}

	for _, table := range tables {
		_, err := s.db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to clean %s: %w", table, err)
		}
	}

	return nil
}