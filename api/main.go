package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/submit-iot-data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			dataSubmissionHandler(w, r)
		} else {
			http.Error(w, METHOD_NOT_SUPPORTED_MSG, http.StatusMethodNotAllowed)
		}
	})

	// Start the web server on port 8080
	fmt.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
