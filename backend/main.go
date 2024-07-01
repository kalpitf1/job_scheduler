package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	job_handlers "github.com/kalpitf1/job_scheduler/backend/job_handlers"
	"github.com/kalpitf1/job_scheduler/backend/websocket"
)

func main() {
	log.Println("Starting Backend")

	r := mux.NewRouter()
	r.HandleFunc("/jobs", job_handlers.GetJobs).Methods("GET")
	r.HandleFunc("/jobs", job_handlers.CreateJob).Methods("POST")
	r.HandleFunc("/ws", websocket.HandleConnections)

	// Apply CORS headers to the router
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Allow requests from your React app
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	go websocket.HandleMessages()

	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}
