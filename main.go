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

	port := ":8080"
	fmt.Printf("Gateaway starting on port%s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}