package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Device represents an IoT device.
type Device struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Submission represents the data submitted by clients.
type Submission struct {
	Timestamp time.Time `json:"timestamp"`
	Location  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Devices []Device `json:"devices"`
}

const METHOD_NOT_SUPPORTED_MSG = "Method not supported"

func dataSubmissionHandler(w http.ResponseWriter, r *http.Request) {
	var submission Submission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		handleError(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	if err := validateSubmission(submission); err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := connectDatabase(w)
	if err != nil {
		return
	}
	defer db.Close()

	if hasAnomalies(submission, db) {
		handleError(w, "An anomaly was detected within the submission", http.StatusBadRequest)
		return
	}

	if err := insertDeviceInteractions(submission, db); err != nil {
		handleError(w, "Error inserting data into the database", http.StatusInternalServerError)
		return
	}

	handleSuccess(w, submission)
}

func handleError(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}

func validateSubmission(submission Submission) error {
	if submission.Timestamp.IsZero() {
		return errors.New("Timestamp is required")
	}

	if submission.Location.Latitude == 0 || submission.Location.Longitude == 0 {
		return errors.New("Latitude and Longitude are required")
	}

	if len(submission.Devices) == 0 {
		return errors.New("At least one device is required")
	}

	return nil
}

func connectDatabase(w http.ResponseWriter) (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/tech_test")
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return nil, err
	}

	return db, nil
}

func insertDeviceInteractions(submission Submission, db *sql.DB) error {
	for _, device := range submission.Devices {
		_, err := db.Exec("INSERT INTO device_interactions (timestamp, latitude, longitude, device_id, device_name) VALUES (?, ?, ?, ?, ?)",
			submission.Timestamp, submission.Location.Latitude, submission.Location.Longitude, device.ID, device.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleSuccess(w http.ResponseWriter, submission Submission) {
	fmt.Printf("Received Data:\nTimestamp: %s\nLocation: %f, %f\nDevices: %v\n",
		submission.Timestamp, submission.Location.Latitude, submission.Location.Longitude, submission.Devices)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data submission successful"))
}

func hasAnomalies(submission Submission, db *sql.DB) bool {
	// Check for anomalies and validate (customize based on your requirements)

	if hasDuplicateTimestamp(submission, db) {
		return true
	}
	return false
}

func hasDuplicateTimestamp(submission Submission, db *sql.DB) bool {

	for _, device := range submission.Devices {
		// Query the db for duplicate timestamps
		query := "SELECT COUNT(*) FROM device_interactions WHERE timestamp = ? AND device_id = ?"
		var count int
		err := db.QueryRow(query, submission.Timestamp, device.ID).Scan(&count)

		if err != nil {
			// Handle db query error
			return true
		}

		if count > 0 {
			// Duplicate timestamp for the same device_id detected
			return true
		}
	}

	return false
}
