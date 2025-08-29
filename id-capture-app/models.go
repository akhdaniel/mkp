package main

import (
	"time"
)

// IDDocument represents the structure of an ID document in the database
type IDDocument struct {
	ID           int       `json:"id"`
	DocumentID   string    `json:"document_id"`
	Name         string    `json:"name"`
	DocumentNumber string  `json:"document_number"`
	BirthDate    string    `json:"birth_date"`
	IssueDate    string    `json:"issue_date"`
	ExpiryDate   string    `json:"expiry_date"`
	ImagePath    string    `json:"image_path"`
	CreatedAt    time.Time `json:"created_at"`
}

// DocumentData represents the data extracted from an ID document
type DocumentData struct {
	Name         string `json:"name"`
	DocumentNumber string `json:"document_number"`
	BirthDate    string `json:"birth_date"`
	IssueDate    string `json:"issue_date"`
	ExpiryDate   string `json:"expiry_date"`
}

// OCRResponse represents the structure of an OCR service response
// This is just an example - the actual structure would depend on the OCR service used
type OCRResponse struct {
	TextAnnotations []TextAnnotation `json:"textAnnotations"`
}

// TextAnnotation represents a text annotation from the OCR service
type TextAnnotation struct {
	Description string `json:"description"`
}