package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// User struct represents the structure of your user data.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	// Add more fields as needed
}

func getUsersFromDatabase(w http.ResponseWriter) ([]User, error) {
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/tech_test")
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	// Query the users table
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and build the users slice
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func exportUsersToCSV(users []User, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"ID", "Name", "Email"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write user data to CSV
	for _, user := range users {
		row := []string{
			fmt.Sprintf("%d", user.ID),
			user.Name,
			user.Email,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	format := r.URL.Query().Get("format")

	users, err := getUsersFromDatabase(w)
	if err != nil {
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}

	var filename string
	var contentType string

	// Determine the response format based on the 'format' query parameter
	switch format {
	case "csv":
		filename = "users.csv"
		contentType = "text/csv"
		err = exportUsersToCSV(users, filename)
		if err != nil {
			http.Error(w, "Error exporting users to CSV", http.StatusInternalServerError)
			return
		}

		// Set response headers for download
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)

		// Open and send the file
		file, err := os.Open(filename)
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Copy the file content to the response writer
		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, "Error sending file", http.StatusInternalServerError)
			return
		}
	case "json":
		filename = "users.json"
		contentType = "application/json"
		response, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(response)
	default:
		http.Error(w, "Invalid or missing 'format' parameter", http.StatusBadRequest)
		return
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to JSON (optional)
	w.Header().Set("Content-Type", "application/json")

	// Create a sample map
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// Create a struct to hold key-value pairs
	type Item struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	// Create a slice to hold items
	var items []Item

	// First for loop using range to iterate through the map
	for key, value := range data {
		items = append(items, Item{Key: key, Value: value})
	}

	// Second for loop using an iterator variable 'i'
	for i := 0; i < len(items); i++ {
		fmt.Printf("Item %d - Key: %s, Value: %s\n", i, items[i].Key, items[i].Value)
	}

	// Construct the JSON response
	response, err := json.Marshal(map[string]interface{}{
		"message": "Hello, World!",
		"items":   items,
	})
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Write(response)
}

func main() {
	// Define a route for the "/hello" endpoint
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			helloWorldHandler(w, r)
		} else {
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsersHandler(w, r)
		} else {
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		}
	})

	// Start the web server on port 8080
	fmt.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
