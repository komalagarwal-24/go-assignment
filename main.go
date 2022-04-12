package main

import (
	"assignment/server"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting HTTP server")
	http.HandleFunc("/getData", server.GetData)

	// Start HTTP server
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
