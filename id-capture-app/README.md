# ID Document Capture Application

This application allows users to upload images of ID documents (National ID or Passport) and extracts key information such as name, document number, birth date, etc. The extracted data is stored in a PostgreSQL database.

## Features

- Upload images of ID documents
- Preview uploaded images before processing
- Extract key information from ID documents (mock implementation)
- Store extracted data in PostgreSQL database
- View document history

## Tech Stack

- **Backend**: Golang with Gin framework
- **Frontend**: React.js
- **Database**: PostgreSQL
- **Deployment**: Docker & Docker Compose

## Prerequisites

- Docker and Docker Compose
- Go (if running without Docker)
- Node.js (if running without Docker)

## Running with Docker (Recommended)

1. Clone the repository
2. Navigate to the project directory
3. Run the application using Docker Compose:

```bash
docker-compose up --build
```

4. Access the application at http://localhost:8080

## Running without Docker

### Backend

1. Install Go dependencies:
```bash
go mod tidy
```

2. Set up PostgreSQL database and update `.env` file with your database credentials

3. Run the backend:
```bash
go run .
```

### Frontend

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm start
```

## API Endpoints

- `POST /api/upload` - Upload an ID document image
- `GET /api/documents` - Get all uploaded documents
- `GET /api/documents/:id` - Get a specific document by ID

## Project Structure

```
id-capture-app/
├── main.go              # Entry point
├── database.go          # Database initialization and connection
├── models.go            # Data structures
├── handlers.go          # HTTP handlers
├── .env                 # Environment variables
├── go.mod               # Go module definition
├── go.sum               # Go dependencies checksums
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
├── uploads/             # Uploaded images (created at runtime)
└── frontend/            # React frontend
    ├── public/
    └── src/
        ├── App.js
        ├── App.css
        ├── index.js
        └── index.css
```

## Implementation Notes

The current implementation includes a mock ID data extraction function. In a production environment, you would replace this with an actual OCR solution or machine learning model for extracting text from ID documents.

## Future Enhancements

- Integrate with a real OCR service (e.g., Google Vision API, Tesseract)
- Add user authentication
- Implement document validation
- Add support for more document types
- Improve UI/UX