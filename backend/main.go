package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Job represents data about a posted job
type Job struct {
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
	Status   string        `json:"status"`
}

var jobs []Job

func main() {
	log.Println("Starting Backend")

	r := mux.NewRouter()
	r.HandleFunc("/jobs", getJobs).Methods("GET")
	r.HandleFunc("/jobs", createJob).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(jobs)
}

func createJob(w http.ResponseWriter, r *http.Request) {
	var newJob Job
	_ = json.NewDecoder(r.Body).Decode(&newJob)
	jobs = append(jobs, newJob)
	json.NewEncoder(w).Encode(newJob)
}
