package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Ensure uploads directory exists
	uploadDir := "../../uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Printf("Warning: could not create uploads dir: %v", err)
	}

	// Serve files
	// The URL path /uploads/filename.jpg should map to ../../uploads/filename.jpg
	fs := http.FileServer(http.Dir(uploadDir))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	port := ":8084"
	fmt.Printf("Media Service running on %s, serving %s\n", port, uploadDir)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
