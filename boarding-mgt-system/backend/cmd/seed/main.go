package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ferryflow/backend/internal/config"
	"github.com/ferryflow/backend/internal/database"
	"github.com/ferryflow/backend/internal/models"
	"github.com/ferryflow/backend/internal/repository"
	"golang.org/x/crypto/argon2"
	"encoding/base64"
	"crypto/rand"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	operatorRepo := repository.NewOperatorRepository(db)
	portRepo := repository.NewPortRepository(db)
	vesselRepo := repository.NewVesselRepository(db)
	routeRepo := repository.NewRouteRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)

	ctx := context.Background()

	fmt.Println("🌱 Starting database seeding...")

	// Create demo users
	users := []struct {
		email    string
		password string
		role     string
		name     string
	}{
		{"admin@demo.com", "password123", "admin", "Admin User"},
		{"operator@demo.com", "password123", "operator", "Operator User"},
		{"customer@demo.com", "password123", "customer", "Customer User"},
		{"john.doe@example.com", "password123", "customer", "John Doe"},
		{"jane.smith@example.com", "password123", "customer", "Jane Smith"},
	}

	userIDs := make(map[string]uuid.UUID)
	for _, u := range users {
		hashedPassword := hashPassword(u.password)
		user := &models.User{
			Email:        u.email,
			PasswordHash: hashedPassword,
			Role:         u.role,
			FullName:     u.name,
			PhoneNumber:  "+1234567890",
			IsVerified:   true,
			IsActive:     true,
		}
		
		err := userRepo.Create(ctx, user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", u.email, err)
			continue
		}
		userIDs[u.email] = user.ID
		fmt.Printf("✅ Created user: %s\n", u.email)
	}

	// Create ferry operators
	operators := []models.Operator{
		{
			Name:         "FastFerry Lines",
			ContactEmail: "contact@fastferry.com",
			ContactPhone: "+1234567890",
			Address:      "123 Harbor St, Port City",
			IsActive:     true,
		},
		{
			Name:         "Island Hoppers",
			ContactEmail: "info@islandhoppers.com",
			ContactPhone: "+0987654321",
			Address:      "456 Marina Blvd, Island Bay",
			IsActive:     true,
		},
	}

	operatorIDs := make([]uuid.UUID, 0)
	for _, op := range operators {
		err := operatorRepo.Create(ctx, &op)
		if err != nil {
			log.Printf("Failed to create operator %s: %v", op.Name, err)
			continue
		}
		operatorIDs = append(operatorIDs, op.ID)
		fmt.Printf("✅ Created operator: %s\n", op.Name)
	}

	// Create ports
	ports := []models.Port{
		{
			Code:        "MNL",
			Name:        "Manila Port",
			City:        "Manila",
			Country:     "Philippines",
			Latitude:    14.5995,
			Longitude:   120.9842,
			Timezone:    "Asia/Manila",
			Facilities:  []string{"Parking", "Restaurant", "Waiting Area", "WiFi"},
			IsActive:    true,
		},
		{
			Code:        "CEB",
			Name:        "Cebu Port",
			City:        "Cebu",
			Country:     "Philippines",
			Latitude:    10.3157,
			Longitude:   123.8854,
			Timezone:    "Asia/Manila",
			Facilities:  []string{"Parking", "Food Court", "Shops", "WiFi"},
			IsActive:    true,
		},
		{
			Code:        "BOH",
			Name:        "Bohol Port",
			City:        "Tagbilaran",
			Country:     "Philippines",
			Latitude:    9.6407,
			Longitude:   123.8543,
			Timezone:    "Asia/Manila",
			Facilities:  []string{"Parking", "Cafeteria", "WiFi"},
			IsActive:    true,
		},
		{
			Code:        "ILO",
			Name:        "Iloilo Port",
			City:        "Iloilo",
			Country:     "Philippines",
			Latitude:    10.6969,
			Longitude:   122.5644,
			Timezone:    "Asia/Manila",
			Facilities:  []string{"Parking", "Restaurant", "Shops"},
			IsActive:    true,
		},
	}

	portIDs := make(map[string]uuid.UUID)
	for _, p := range ports {
		err := portRepo.Create(ctx, &p)
		if err != nil {
			log.Printf("Failed to create port %s: %v", p.Name, err)
			continue
		}
		portIDs[p.Code] = p.ID
		fmt.Printf("✅ Created port: %s (%s)\n", p.Name, p.Code)
	}

	// Create vessels
	vessels := []models.Vessel{
		{
			OperatorID:       operatorIDs[0],
			Name:            "MV Sea Explorer",
			RegistrationNo:  "REG001",
			Capacity:        300,
			VesselType:      "passenger",
			YearBuilt:       2020,
			Length:          85.5,
			Width:           15.2,
			MaxSpeed:        25.5,
			Features:        []string{"Air Conditioning", "Cafeteria", "WiFi", "Entertainment"},
			IsActive:        true,
		},
		{
			OperatorID:       operatorIDs[0],
			Name:            "MV Ocean Pride",
			RegistrationNo:  "REG002",
			Capacity:        250,
			VesselType:      "passenger",
			YearBuilt:       2019,
			Length:          75.0,
			Width:           14.0,
			MaxSpeed:        23.0,
			Features:        []string{"Air Conditioning", "Restaurant", "WiFi"},
			IsActive:        true,
		},
		{
			OperatorID:       operatorIDs[1],
			Name:            "MV Island Express",
			RegistrationNo:  "REG003",
			Capacity:        200,
			VesselType:      "passenger",
			YearBuilt:       2021,
			Length:          70.0,
			Width:           12.5,
			MaxSpeed:        28.0,
			Features:        []string{"Air Conditioning", "Snack Bar", "WiFi", "Deck Seating"},
			IsActive:        true,
		},
	}

	vesselIDs := make([]uuid.UUID, 0)
	for _, v := range vessels {
		err := vesselRepo.Create(ctx, &v)
		if err != nil {
			log.Printf("Failed to create vessel %s: %v", v.Name, err)
			continue
		}
		vesselIDs = append(vesselIDs, v.ID)
		fmt.Printf("✅ Created vessel: %s\n", v.Name)
	}

	// Create routes
	routes := []models.Route{
		{
			OperatorID:       operatorIDs[0],
			DeparturePortID:  portIDs["MNL"],
			ArrivalPortID:    portIDs["CEB"],
			Distance:         572.0,
			EstimatedDuration: 22 * 60, // 22 hours
			BasePrice:        2500.00,
			IsActive:         true,
		},
		{
			OperatorID:       operatorIDs[0],
			DeparturePortID:  portIDs["CEB"],
			ArrivalPortID:    portIDs["BOH"],
			Distance:         72.0,
			EstimatedDuration: 2 * 60, // 2 hours
			BasePrice:        500.00,
			IsActive:         true,
		},
		{
			OperatorID:       operatorIDs[1],
			DeparturePortID:  portIDs["CEB"],
			ArrivalPortID:    portIDs["ILO"],
			Distance:         186.0,
			EstimatedDuration: 6 * 60, // 6 hours
			BasePrice:        1200.00,
			IsActive:         true,
		},
		{
			OperatorID:       operatorIDs[1],
			DeparturePortID:  portIDs["ILO"],
			ArrivalPortID:    portIDs["MNL"],
			Distance:         460.0,
			EstimatedDuration: 18 * 60, // 18 hours
			BasePrice:        2000.00,
			IsActive:         true,
		},
	}

	routeIDs := make([]uuid.UUID, 0)
	for _, r := range routes {
		err := routeRepo.Create(ctx, &r)
		if err != nil {
			log.Printf("Failed to create route: %v", err)
			continue
		}
		routeIDs = append(routeIDs, r.ID)
		fmt.Printf("✅ Created route: %s to %s\n", r.DeparturePortID, r.ArrivalPortID)
	}

	// Create schedules for the next 30 days
	now := time.Now()
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, i)
		
		// Morning schedule: Manila to Cebu
		if i%2 == 0 { // Every other day
			schedule1 := models.Schedule{
				RouteID:        routeIDs[0],
				VesselID:       vesselIDs[0],
				DepartureTime:  time.Date(date.Year(), date.Month(), date.Day(), 6, 0, 0, 0, time.UTC),
				ArrivalTime:    time.Date(date.Year(), date.Month(), date.Day()+1, 4, 0, 0, 0, time.UTC),
				Price:          2500.00,
				AvailableSeats: 300,
				TotalSeats:     300,
				Status:         "scheduled",
			}
			err := scheduleRepo.Create(ctx, &schedule1)
			if err != nil {
				log.Printf("Failed to create schedule: %v", err)
			}
		}
		
		// Daily schedule: Cebu to Bohol
		schedule2 := models.Schedule{
			RouteID:        routeIDs[1],
			VesselID:       vesselIDs[1],
			DepartureTime:  time.Date(date.Year(), date.Month(), date.Day(), 8, 0, 0, 0, time.UTC),
			ArrivalTime:    time.Date(date.Year(), date.Month(), date.Day(), 10, 0, 0, 0, time.UTC),
			Price:          500.00,
			AvailableSeats: 250,
			TotalSeats:     250,
			Status:         "scheduled",
		}
		err := scheduleRepo.Create(ctx, &schedule2)
		if err != nil {
			log.Printf("Failed to create schedule: %v", err)
		}
		
		// Afternoon schedule: Cebu to Bohol
		schedule3 := models.Schedule{
			RouteID:        routeIDs[1],
			VesselID:       vesselIDs[2],
			DepartureTime:  time.Date(date.Year(), date.Month(), date.Day(), 14, 0, 0, 0, time.UTC),
			ArrivalTime:    time.Date(date.Year(), date.Month(), date.Day(), 16, 0, 0, 0, time.UTC),
			Price:          500.00,
			AvailableSeats: 200,
			TotalSeats:     200,
			Status:         "scheduled",
		}
		err := scheduleRepo.Create(ctx, &schedule3)
		if err != nil {
			log.Printf("Failed to create schedule: %v", err)
		}
	}
	
	fmt.Println("✅ Created schedules for the next 30 days")

	fmt.Println("\n🎉 Database seeding completed successfully!")
	fmt.Println("\nDemo accounts created:")
	fmt.Println("  Admin:    admin@demo.com / password123")
	fmt.Println("  Operator: operator@demo.com / password123")
	fmt.Println("  Customer: customer@demo.com / password123")
}

func hashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, 64*1024, 1, 4, b64Salt, b64Hash)
}