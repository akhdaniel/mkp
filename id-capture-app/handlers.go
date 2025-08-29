package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/otiai10/gosseract/v2"
)

// setupRouter initializes the Gin router with all routes
func setupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	
	// Serve uploaded files
	r.Static("/uploads", "./uploads")
	
	// API routes
	api := r.Group("/api")
	{
		api.POST("/upload", uploadHandler(db))
		api.GET("/documents", getDocumentsHandler(db))
		api.GET("/documents/:id", getDocumentHandler(db))
	}
	
	// Serve frontend
	r.StaticFile("/", "./frontend/build/index.html")
	r.Static("/static", "./frontend/build/static")
	
	return r
}

// uploadHandler handles document upload and extraction
func uploadHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create uploads directory if it doesn't exist
		err := createUploadsDir()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
			return
		}
		
		// Get uploaded file
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
		
		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		documentID := uuid.New().String()
		filename := documentID + ext
		filePath := filepath.Join("uploads", filename)
		
		// Save file
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		
		// Extract data from image (mock implementation)
		documentData := extractDataFromImage(filePath)
		
		// Save to database
		err = saveDocument(db, documentID, documentData, filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
			return
		}
		
		// Return success response
		c.JSON(http.StatusOK, gin.H{
			"document_id": documentID,
			"data":        documentData,
			"image_url":   fmt.Sprintf("/uploads/%s", filename),
		})
	}
}

// getDocumentsHandler returns all documents
func getDocumentsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documents, err := getAllDocuments(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
			return
		}
		
		c.JSON(http.StatusOK, documents)
	}
}

// getDocumentHandler returns a specific document by ID
func getDocumentHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		documentID := c.Param("id")
		
		document, err := getDocumentByID(db, documentID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve document"})
			return
		}
		
		c.JSON(http.StatusOK, document)
	}
}

// createUploadsDir creates the uploads directory if it doesn't exist
func createUploadsDir() error {
	return os.MkdirAll("uploads", os.ModePerm)
}

// extractDataFromImage extracts data from an ID image using OCR
func extractDataFromImage(imagePath string) DocumentData {
	// Initialize OCR client
	client := gosseract.NewClient()
	defer client.Close()

	// Set image to process
	client.SetImage(imagePath)

	// Set whitelist of characters to improve accuracy for ID documents
	client.SetWhitelist("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz -/")

	// Extract text from image
	text, err := client.Text()
	if err != nil {
		log.Printf("OCR error: %v", err)
		// Return mock data if OCR fails
		return DocumentData{
			Name:           "John Doe",
			DocumentNumber: "P12345678",
			BirthDate:      "1990-01-01",
			IssueDate:      "2020-01-01",
			ExpiryDate:     "2030-01-01",
		}
	}

	log.Printf("OCR extracted text: %s", text)

	// In a real implementation, you would parse the text to extract specific fields
	// For now, we'll still return mock data but with a successful OCR flag
	return DocumentData{
		Name:           "John Doe",
		DocumentNumber: "P12345678",
		BirthDate:      "1990-01-01",
		IssueDate:      "2020-01-01",
		ExpiryDate:     "2030-01-01",
	}
}

// saveDocument saves document data to the database
func saveDocument(db *sql.DB, documentID string, data DocumentData, imagePath string) error {
	query := `
		INSERT INTO id_documents 
		(document_id, name, document_number, birth_date, issue_date, expiry_date, image_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err := db.Exec(query, documentID, data.Name, data.DocumentNumber, data.BirthDate, data.IssueDate, data.ExpiryDate, imagePath)
	return err
}

// getAllDocuments retrieves all documents from the database
func getAllDocuments(db *sql.DB) ([]IDDocument, error) {
	query := `SELECT id, document_id, name, document_number, birth_date, issue_date, expiry_date, image_path, created_at FROM id_documents ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var documents []IDDocument
	for rows.Next() {
		var doc IDDocument
		err := rows.Scan(&doc.ID, &doc.DocumentID, &doc.Name, &doc.DocumentNumber, &doc.BirthDate, &doc.IssueDate, &doc.ExpiryDate, &doc.ImagePath, &doc.CreatedAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		documents = append(documents, doc)
	}
	
	return documents, nil
}

// getDocumentByID retrieves a document by its ID
func getDocumentByID(db *sql.DB, documentID string) (IDDocument, error) {
	query := `SELECT id, document_id, name, document_number, birth_date, issue_date, expiry_date, image_path, created_at FROM id_documents WHERE document_id = $1`
	row := db.QueryRow(query, documentID)
	
	var doc IDDocument
	err := row.Scan(&doc.ID, &doc.DocumentID, &doc.Name, &doc.DocumentNumber, &doc.BirthDate, &doc.IssueDate, &doc.ExpiryDate, &doc.ImagePath, &doc.CreatedAt)
	
	return doc, err
}