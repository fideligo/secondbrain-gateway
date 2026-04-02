package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type DocumentRequest struct {
	FileName string `json:"file_name"`
	Author string `json:"author"`
}

func main() {
	
	// Health Check (GET)
	http.HandleFunc("/api/health", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "SecondBrain Gateaway is running smoothly!")
	})

	// Menerima data JSON (POST)
	http.HandleFunc("/api/upload", func(writer http.ResponseWriter, request *http.Request) {

		if request.Method != http.MethodPost {
			http.Error(writer, "Only accepts POST method", http.StatusMethodNotAllowed)
			return
		}

		var docReq DocumentRequest

		err := json.NewDecoder(request.Body).Decode(&docReq)
		if err != nil {
			http.Error(writer, "Invalid JSON Format", http.StatusBadRequest)
			return
		}

		fmt.Printf("New Document!\nFile Name: %s\nWriter: %s\n", docReq.FileName, docReq.Author)

		writer.WriteHeader(http.StatusAccepted) // Memberikan kode status HTTP 202 Accepted
		fmt.Fprintf(writer, "The Document %s from %s is accepted and being processed.", docReq.FileName, docReq.Author)
	})

	port := ":8080"
	fmt.Printf("Gateaway starting on port%s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}