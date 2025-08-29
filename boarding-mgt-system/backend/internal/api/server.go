package api

import (
	"time"

	"github.com/ferryflow/boarding-mgt-system/internal/api/handlers"
	"github.com/ferryflow/boarding-mgt-system/internal/api/middleware"
	"github.com/ferryflow/boarding-mgt-system/internal/config"
	"github.com/ferryflow/boarding-mgt-system/internal/database"
	"github.com/ferryflow/boarding-mgt-system/internal/repository"
	"github.com/ferryflow/boarding-mgt-system/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Router   *gin.Engine
	config   *config.Config
	db       *database.DB
	services *service.Services
}

func NewServer(cfg *config.Config, db *database.DB) *Server {
	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	
	// Initialize repositories
	repos := repository.NewRepositories(db)
	
	// Initialize services
	services := service.NewServices(repos, cfg)
	
	server := &Server{
		Router:   router,
		config:   cfg,
		db:       db,
		services: services,
	}
	
	server.setupMiddleware()
	server.setupRoutes()
	
	return server
}

func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.Router.Use(gin.Recovery())
	
	// Logger middleware
	s.Router.Use(gin.Logger())
	
	// CORS middleware
	s.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	// Request ID middleware
	s.Router.Use(middleware.RequestID())
	
	// Rate limiting middleware
	s.Router.Use(middleware.RateLimit())
}

func (s *Server) setupRoutes() {
	// Health check
	s.Router.GET("/health", s.healthCheck)
	
	// Swagger documentation
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// API v1 routes
	v1 := s.Router.Group("/api/v1")
	
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(s.services.Auth)
	operatorHandler := handlers.NewOperatorHandler(s.services.Operator)
	portHandler := handlers.NewPortHandler(s.services.Port)
	vesselHandler := handlers.NewVesselHandler(s.services.Vessel)
	routeHandler := handlers.NewRouteHandler(s.services.Route)
	scheduleHandler := handlers.NewScheduleHandler(s.services.Schedule)
	bookingHandler := handlers.NewBookingHandler(s.services.Booking)
	ticketHandler := handlers.NewTicketHandler(s.services.Ticket)
	userHandler := handlers.NewUserHandler(s.services.User)
	
	// Public routes (no authentication required)
	public := v1.Group("")
	{
		// Authentication
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", authHandler.RefreshToken)
		
		// Public schedule search
		public.GET("/schedules/search", scheduleHandler.SearchSchedules)
		public.GET("/schedules/:id", scheduleHandler.GetSchedule)
		
		// Public port information
		public.GET("/ports", portHandler.ListPorts)
		public.GET("/ports/:id", portHandler.GetPort)
		
		// Public route information
		public.GET("/routes", routeHandler.ListRoutes)
		public.GET("/routes/:id", routeHandler.GetRoute)
	}
	
	// Protected routes (authentication required)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(s.config.JWT))
	{
		// User profile
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.PUT("/users/profile", userHandler.UpdateProfile)
		protected.POST("/auth/logout", authHandler.Logout)
		
		// Customer bookings
		protected.GET("/bookings/my", bookingHandler.GetMyBookings)
		protected.POST("/bookings", bookingHandler.CreateBooking)
		protected.GET("/bookings/:id", bookingHandler.GetBooking)
		protected.POST("/bookings/:id/cancel", bookingHandler.CancelBooking)
		
		// Tickets
		protected.GET("/tickets/my", ticketHandler.GetMyTickets)
		protected.GET("/tickets/:id", ticketHandler.GetTicket)
		protected.GET("/tickets/:id/qr", ticketHandler.GetTicketQR)
	}
	
	// Admin routes (operator/admin authentication required)
	admin := v1.Group("")
	admin.Use(middleware.AuthMiddleware(s.config.JWT))
	admin.Use(middleware.RequireRole("agent", "operator_admin", "system_admin"))
	{
		// Operator management
		admin.GET("/operators", operatorHandler.ListOperators)
		admin.GET("/operators/:id", operatorHandler.GetOperator)
		admin.POST("/operators", middleware.RequireRole("system_admin"), operatorHandler.CreateOperator)
		admin.PUT("/operators/:id", operatorHandler.UpdateOperator)
		admin.DELETE("/operators/:id", middleware.RequireRole("system_admin"), operatorHandler.DeleteOperator)
		
		// Port management
		admin.POST("/ports", portHandler.CreatePort)
		admin.PUT("/ports/:id", portHandler.UpdatePort)
		admin.DELETE("/ports/:id", portHandler.DeletePort)
		
		// Vessel management
		admin.GET("/vessels", vesselHandler.ListVessels)
		admin.GET("/vessels/:id", vesselHandler.GetVessel)
		admin.POST("/vessels", vesselHandler.CreateVessel)
		admin.PUT("/vessels/:id", vesselHandler.UpdateVessel)
		admin.DELETE("/vessels/:id", vesselHandler.DeleteVessel)
		
		// Route management
		admin.POST("/routes", routeHandler.CreateRoute)
		admin.PUT("/routes/:id", routeHandler.UpdateRoute)
		admin.DELETE("/routes/:id", routeHandler.DeleteRoute)
		
		// Schedule management
		admin.GET("/schedules", scheduleHandler.ListSchedules)
		admin.POST("/schedules", scheduleHandler.CreateSchedule)
		admin.PUT("/schedules/:id", scheduleHandler.UpdateSchedule)
		admin.DELETE("/schedules/:id", scheduleHandler.DeleteSchedule)
		admin.POST("/schedules/:id/cancel", scheduleHandler.CancelSchedule)
		
		// Booking management
		admin.GET("/bookings", bookingHandler.ListBookings)
		admin.PUT("/bookings/:id", bookingHandler.UpdateBooking)
		
		// User management
		admin.GET("/users", middleware.RequireRole("operator_admin", "system_admin"), userHandler.ListUsers)
		admin.GET("/users/:id", userHandler.GetUser)
		admin.PUT("/users/:id", middleware.RequireRole("operator_admin", "system_admin"), userHandler.UpdateUser)
		admin.DELETE("/users/:id", middleware.RequireRole("system_admin"), userHandler.DeleteUser)
		
		// Reports
		admin.GET("/reports/bookings", bookingHandler.GetBookingReport)
		admin.GET("/reports/revenue", bookingHandler.GetRevenueReport)
		admin.GET("/reports/manifest/:schedule_id", scheduleHandler.GetManifest)
	}
	
	// System admin only routes
	systemAdmin := v1.Group("")
	systemAdmin.Use(middleware.AuthMiddleware(s.config.JWT))
	systemAdmin.Use(middleware.RequireRole("system_admin"))
	{
		// Audit logs
		systemAdmin.GET("/audit-logs", s.getAuditLogs)
		
		// System stats
		systemAdmin.GET("/system/stats", s.getSystemStats)
	}
}

// Health check endpoint
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"time":   time.Now().Unix(),
	})
}

// Placeholder for audit logs
func (s *Server) getAuditLogs(c *gin.Context) {
	// TODO: Implement audit log retrieval
	c.JSON(200, gin.H{
		"message": "Audit logs endpoint",
	})
}

// Placeholder for system stats
func (s *Server) getSystemStats(c *gin.Context) {
	// TODO: Implement system statistics
	c.JSON(200, gin.H{
		"message": "System stats endpoint",
	})
}