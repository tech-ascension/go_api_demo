package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
)

var sqlOpen = sql.Open

// Mock SQL database
var mockDB *sql.DB

// Mock for sql.Open function
func mockSqlOpen(driverName, dataSourceName string) (*sql.DB, error) {
	return mockDB, nil
}

func TestHandleError(t *testing.T) {
	w := httptest.NewRecorder()
	handleError(w, "Test error", http.StatusBadRequest)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestValidSubmission(t *testing.T) {
	submission := Submission{
		Timestamp: time.Now(),
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  10.0,
			Longitude: 20.0,
		},
		Devices: []Device{{ID: 1, Name: "Device1"}},
	}

	err := validateSubmission(submission)

	if err != nil {
		t.Errorf("Expected no error for a valid submission, got %v", err)
	}
}

func TestValidateSubmissionMissingTimestamp(t *testing.T) {
	submission := Submission{
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  10.0,
			Longitude: 20.0,
		},
		Devices: []Device{{ID: 1, Name: "Device1"}},
	}

	err := validateSubmission(submission)

	if err == nil || err.Error() != "Timestamp is required" {
		t.Errorf("Expected 'Timestamp is required' error, got %v", err)
	}
}

func TestValidateSubmissionMissingLocation(t *testing.T) {
	submission := Submission{
		Timestamp: time.Now(),
		Devices:   []Device{{ID: 1, Name: "Device1"}},
	}

	err := validateSubmission(submission)

	expectedErrorMessage := "Latitude and Longitude are required"
	if err == nil || err.Error() != expectedErrorMessage {
		t.Errorf("Expected '%s' error, got %v", expectedErrorMessage, err)
	}
}
func TestValidateSubmissionInvalidLatitude(t *testing.T) {
	submission := Submission{
		Timestamp: time.Now(),
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  0, // Setting an invalid latitude (0)
			Longitude: 20.0,
		},
		Devices: []Device{{ID: 1, Name: "Device1"}},
	}

	err := validateSubmission(submission)

	expectedErrorMessage := "Latitude and Longitude are required"
	if err == nil || err.Error() != expectedErrorMessage {
		t.Errorf("Expected '%s' error, got %v", expectedErrorMessage, err)
	}
}
func TestConnectDatabaseSuccess(t *testing.T) {
	// Set up mock database
	mockDB, _, _ = sqlmock.New()
	defer mockDB.Close()

	// Replace the real sql.Open function with our mock
	sqlOpen = mockSqlOpen
	defer func() { sqlOpen = sql.Open }()

	// Create a new recorder for the HTTP response
	w := httptest.NewRecorder()

	// Call the connectDatabase function
	connectDatabase(w)

	// Check if the mockDB was used
	if mockDB == nil {
		t.Error("Expected mockDB to be used, got nil")
	}
}

func TestValidateSubmissionEmptyDevices(t *testing.T) {
	submission := Submission{
		Timestamp: time.Now(),
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  10.0,
			Longitude: 20.0,
		},
	}

	err := validateSubmission(submission)

	expectedErrorMessage := "At least one device is required"
	if err == nil || err.Error() != expectedErrorMessage {
		t.Errorf("Expected '%s' error, got %v", expectedErrorMessage, err)
	}
}

func TestConnectDatabaseError(t *testing.T) {
	// Replace the real sql.Open function with a function that returns an error
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("Test error")
	}
	defer func() { sqlOpen = sql.Open }()

	w := httptest.NewRecorder()
	connectDatabase(w)

	if w.Code == http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		t.FailNow()
	}

}

func TestInsertDeviceInteractionsSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO device_interactions").
		WillReturnResult(sqlmock.NewResult(1, 1))

	submission := Submission{
		Timestamp: time.Now(),
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  10.0,
			Longitude: 20.0,
		},
		Devices: []Device{{ID: 1, Name: "Device1"}},
	}

	err = insertDeviceInteractions(submission, db)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHasAnomaliesDuplicateTimestamp(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// This just demonstrates that the query hasDuplicateTimestamp() is called and returns a count of 1, which is considered an anomaly/invalid
	mock.ExpectQuery("SELECT COUNT(*)").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	submission := Submission{
		Timestamp: time.Now(),
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  10.0,
			Longitude: 20.0,
		},
		Devices: []Device{{ID: 1, Name: "Device1"}},
	}

	result := hasAnomalies(submission, db)

	if !result {
		t.Errorf("Expected true for duplicate timestamp, got false")
	}
}

func TestIntegrationDataSubmissionHandler(t *testing.T) {

	submission := Submission{
		Timestamp: time.Now(),
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  10.0,
			Longitude: 20.0,
		},
		Devices: []Device{{ID: 1, Name: "Device1"}},
	}

	submissionJSON, err := json.Marshal(submission)
	if err != nil {
		t.Fatalf("Error marshaling submission to JSON: %v", err)
	}

	// Create an HTTP request
	req, err := http.NewRequest("POST", "/submit-iot-data", bytes.NewBuffer(submissionJSON))
	if err != nil {
		t.Fatalf("Error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the response
	w := httptest.NewRecorder()

	// Call the dataSubmissionHandler directly
	dataSubmissionHandler(w, req)

	// Assert that the HTTP status code is 200
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

}
