package job_handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	job_handlers "github.com/kalpitf1/job_scheduler/backend/job_handlers"
	"github.com/kalpitf1/job_scheduler/backend/models"
	"github.com/kalpitf1/job_scheduler/backend/websocket"
)

// Create a new mux router and register the GetJobs handler.
func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/jobs", job_handlers.GetJobs).Methods("GET")
	r.HandleFunc("/jobs", job_handlers.CreateJob).Methods("POST")
	r.HandleFunc("/ws", websocket.HandleConnections)
	return r
}

// Setup: Initialize the jobs slice with some test data.
func addTestJobs(numOfJobs int) {
	job_handlers.JobsMutex.Lock()
	defer job_handlers.JobsMutex.Unlock()
	job_handlers.Jobs = nil
	for i := 0; i < numOfJobs; i++ {
		job_handlers.Jobs = append(job_handlers.Jobs, &models.Job{
			Name:     "Job " + string(rune(i+1)),
			Duration: time.Duration(i+1) * time.Second,
			Status:   "Completed",
			ID:       i + 1,
		})
	}
}

func TestGetJobs(t *testing.T) {
	// Initialize the handler package's state
	job_handlers.Reset()

	// Setup: Initialize the jobs slice with 2 jobs.
	numOfJobs := 2
	addTestJobs(numOfJobs)

	// Create a new HTTP request to the /jobs endpoint.
	req, err := http.NewRequest("GET", "/jobs", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response.
	w := httptest.NewRecorder()

	// Setup the router for API calls
	r := setupRouter()

	// Serve the HTTP request.
	r.ServeHTTP(w, req)

	// Check the status code.
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body.
	var responseJobs []*models.Job
	err = json.NewDecoder(w.Body).Decode(&responseJobs)
	if err != nil {
		t.Errorf("error decoding response: %v", err)
	}

	if len(responseJobs) != numOfJobs {
		t.Errorf("handler returned unexpected number of jobs: got %v want %v", len(responseJobs), 2)
	}

	expected := job_handlers.Jobs
	for i, job := range responseJobs {
		if job.Name != expected[i].Name || job.Duration != expected[i].Duration || job.Status != expected[i].Status || job.ID != expected[i].ID {
			t.Errorf("handler returned unexpected job: got %+v want %+v", job, expected[i])
		}
	}
}

func TestCreateJob(t *testing.T) {
	// Initialize the handler package's state
	job_handlers.Reset()

	// Disable broadcasting for the test
	job_handlers.BroadcastEnabledMutex.Lock()
	job_handlers.BroadcastEnabled = false
	job_handlers.BroadcastEnabledMutex.Unlock()
	defer func() {
		job_handlers.BroadcastEnabledMutex.Lock()
		job_handlers.BroadcastEnabled = true
		job_handlers.BroadcastEnabledMutex.Unlock()

	}() // Ensure it's re-enabled after the test

	done := make(chan bool)
	go func() {
		newJob := models.Job{
			Name:     "Job 1",
			Duration: time.Duration(1) * time.Second,
		}

		body, err := json.Marshal(newJob)
		if err != nil {
			t.Fatalf("Failed to marshal job: %v", err)
		}

		// Create a new HTTP request to the /jobs endpoint.
		req, err := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to capture the response.
		w := httptest.NewRecorder()

		// Setup the router for API calls
		r := setupRouter()

		// Serve the HTTP request.
		// log.Println("TestCreateJob: Calling ServeHTTP")
		r.ServeHTTP(w, req)
		// log.Println("TestCreateJob: ServeHTTP returned")

		// Check the status code.
		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		done <- true
	}()

	select {
	case <-done:
	// Test completed successfully
	case <-time.After(10 * time.Second):
		t.Fatal("TestCreateJob timed out")
	}
}

func BenchmarkGetJobs(b *testing.B) {
	// Add numOfJobs jobs for the benchmark
	numOfJobs := 1000
	addTestJobs(numOfJobs)

	// Create a new HTTP request to the /jobs endpoint.
	req, err := http.NewRequest("GET", "/jobs", nil)
	if err != nil {
		b.Fatal(err)
	}

	// Setup the router for API calls
	r := setupRouter()

	// Reset timer due to expensive setup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a ResponseRecorder to capture the response.
		w := httptest.NewRecorder()

		// Serve the HTTP request.
		r.ServeHTTP(w, req)

		// Check the status code.
		if status := w.Code; status != http.StatusOK {
			b.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Check the response body.
		var responseJobs []*models.Job
		err = json.NewDecoder(w.Body).Decode(&responseJobs)
		if err != nil {
			b.Errorf("error decoding response: %v", err)
		}

		if len(responseJobs) != numOfJobs {
			b.Errorf("handler returned unexpected number of jobs: got %v want %v", len(responseJobs), 2)
		}
	}
}

func BenchmarkCreateJob(b *testing.B) {
	// Initialize the handler package's state
	job_handlers.Reset()

	// Disable broadcasting for the test
	job_handlers.BroadcastEnabledMutex.Lock()
	job_handlers.BroadcastEnabled = false
	job_handlers.BroadcastEnabledMutex.Unlock()

	defer func() {
		job_handlers.BroadcastEnabledMutex.Lock()
		job_handlers.BroadcastEnabled = true
		job_handlers.BroadcastEnabledMutex.Unlock()

	}() // Ensure it's re-enabled after the test

	newJob := models.Job{
		Name:     "Job 1",
		Duration: time.Duration(1) * time.Second,
	}

	body, err := json.Marshal(newJob)
	if err != nil {
		b.Fatalf("Failed to marshal job: %v", err)
	}

	// Create a new HTTP request to the /jobs endpoint.
	// req, err := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	// if err != nil {
	// 	b.Fatal(err)
	// }

	// Setup the router for API calls
	r := setupRouter()

	// Reset timer due to expensive setup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		job_handlers.Reset()
		job_handlers.BroadcastEnabledMutex.Lock()
		job_handlers.BroadcastEnabled = false
		job_handlers.BroadcastEnabledMutex.Unlock()

		req, err := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
		if err != nil {
			b.Fatal(err)
		}

		// Create a ResponseRecorder to capture the response.
		w := httptest.NewRecorder()

		// Serve the HTTP request.
		r.ServeHTTP(w, req)

		// Check the status code.
		if status := w.Code; status != http.StatusOK {
			b.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}
}
