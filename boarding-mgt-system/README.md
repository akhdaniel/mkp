# FerryFlow - Ferry Boarding Management System

A comprehensive boarding management system for ferry operations featuring ticket purchasing, scheduling, seat management, and customer service tools.

## 🚀 Quick Start Demo

### Prerequisites
- Docker and Docker Compose installed
- Node.js 20+ (for local development)
- Go 1.22+ (for local development)

### Running with Docker (Recommended for Demo)

1. Clone the repository:
```bash
git clone <repository-url>
cd boarding-mgt-system
```

2. Start the PostgreSQL database:
```bash
docker-compose up -d postgres
```

3. Run database migrations:
```bash
cd backend
go run cmd/migrate/main.go up
```

4. Start the backend server:
```bash
go run cmd/server/main.go
```

5. In a new terminal, start the frontend:
```bash
cd frontend
npm install
npm run dev
```

6. Access the application:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080/api

### Demo Credentials

The system includes demo accounts for testing:

- **Customer Account:**
  - Email: customer@demo.com
  - Password: password123

- **Operator Account:**
  - Email: operator@demo.com
  - Password: password123

- **Admin Account:**
  - Email: admin@demo.com
  - Password: password123

## 🎯 Features

### Core Functionality
- **Online Ticket Purchasing**: 24/7 booking with real-time availability
- **Dynamic Scheduling**: Manage complex multi-route schedules
- **Automated Seat Management**: Smart allocation based on preferences
- **Digital Tickets**: QR code-based boarding passes
- **Payment Processing**: Multiple payment methods support
- **Booking Management**: Easy rescheduling and refunds

### User Types
- **Customers**: Search schedules, book tickets, manage bookings
- **Operators**: Manage vessels, routes, and schedules
- **Administrators**: Full system access and reporting

## 🛠 Technology Stack

### Backend
- **Language**: Go (Golang) 1.22+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 17
- **Authentication**: JWT with refresh tokens
- **Password Hashing**: Argon2id
- **Session Store**: Redis (optional)

### Frontend
- **Framework**: React 18 with Vite
- **Styling**: Tailwind CSS
- **UI Components**: Headless UI, Heroicons
- **State Management**: React Context API
- **API Client**: Axios
- **Date Handling**: date-fns

## 📁 Project Structure

```
boarding-mgt-system/
├── backend/
│   ├── cmd/
│   │   ├── server/        # API server entry point
│   │   └── migrate/       # Database migration tool
│   ├── internal/
│   │   ├── api/          # HTTP handlers and middleware
│   │   ├── auth/         # Authentication logic
│   │   ├── config/       # Configuration management
│   │   ├── database/     # Database connection and migrations
│   │   ├── models/       # Data models
│   │   ├── repository/   # Data access layer
│   │   └── service/      # Business logic layer
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── api/          # API client modules
│   │   ├── components/   # Reusable React components
│   │   ├── contexts/     # React contexts
│   │   └── pages/        # Page components
│   └── package.json
└── docker-compose.yml
```

## 🔧 Development Setup

### Backend Development

1. Set up environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=ferryflow
export DB_PASSWORD=ferryflow_dev_2024
export DB_NAME=ferryflow_dev
export JWT_SECRET=your-secret-key
```

2. Install dependencies:
```bash
cd backend
go mod download
```

3. Run migrations:
```bash
go run cmd/migrate/main.go up
```

4. Start the server:
```bash
go run cmd/server/main.go
```

### Frontend Development

1. Install dependencies:
```bash
cd frontend
npm install
```

2. Create `.env` file:
```env
VITE_API_URL=http://localhost:8080/api
```

3. Start development server:
```bash
npm run dev
```

## 📊 Database Schema

The system uses a comprehensive database schema including:
- Users and authentication
- Ferry operators and ports
- Vessels and routes
- Schedules and bookings
- Tickets and payments
- Support tickets and audit logs

## 🔒 Security Features

- **Password Security**: Argon2id hashing
- **JWT Authentication**: Access and refresh tokens
- **Session Management**: Secure session handling
- **Rate Limiting**: API rate limiting per user
- **CORS Configuration**: Proper CORS setup
- **SQL Injection Prevention**: Parameterized queries

## 📝 API Documentation

The API follows RESTful conventions:

### Authentication Endpoints
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - User logout
- `GET /api/auth/profile` - Get user profile
- `PUT /api/auth/profile` - Update profile

### Booking Endpoints
- `GET /api/schedules/search` - Search available schedules
- `POST /api/bookings` - Create new booking
- `GET /api/bookings` - List user bookings
- `GET /api/bookings/:id` - Get booking details
- `POST /api/bookings/:id/cancel` - Cancel booking

### Admin Endpoints
- `GET /api/operators` - List operators
- `GET /api/ports` - List ports
- `GET /api/vessels` - List vessels
- `GET /api/routes` - List routes

## 🚢 Demo Workflow

1. **Register/Login**: Create an account or use demo credentials
2. **Search Schedules**: Select departure/arrival ports and date
3. **Select Schedule**: Choose from available ferry schedules
4. **Enter Passenger Details**: Add passenger information
5. **Payment**: Select payment method
6. **Confirmation**: Receive booking confirmation with QR code tickets
7. **Manage Bookings**: View, reschedule, or cancel bookings

## 🤝 Contributing

This is a demo application for showcasing a ferry booking system MVP. For production use, additional features and security hardening would be required.

## 📄 License

This project is for demonstration purposes.

---

**Note**: This is an MVP demonstration. Production deployment would require additional security measures, payment gateway integration, and scalability considerations.