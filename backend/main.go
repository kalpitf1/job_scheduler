package main

import (
	"container/heap"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Job represents data about a posted job
type Job struct {
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
	Status   string        `json:"status"`
	Index    int           `json:"-"` // Index in the heap
}

// JobPriorityQueue implements heap.Interface and holds Jobs
type JobPriorityQueue []*Job

func (pq JobPriorityQueue) Len() int { return len(pq) }

func (pq JobPriorityQueue) Less(i, j int) bool {
	return pq[i].Duration < pq[j].Duration
}

func (pq JobPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *JobPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	job := x.(*Job)
	job.Index = n
	*pq = append(*pq, job)
}

func (pq *JobPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	job := old[n-1]
	job.Index = -1 // for safety
	*pq = old[0 : n-1]
	return job
}

func (pq *JobPriorityQueue) Peek() *Job {
	if len(*pq) == 0 {
		return nil
	}
	return (*pq)[0]
}

var (
	jobs        JobPriorityQueue
	jobsMutex   sync.Mutex
	processing  bool
	processingM sync.Mutex
)

func main() {
	log.Println("Starting Backend")

	heap.Init(&jobs)

	r := mux.NewRouter()
	r.HandleFunc("/jobs", getJobs).Methods("GET")
	r.HandleFunc("/jobs", createJob).Methods("POST")

	// Apply CORS headers to the router
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Allow requests from your React app
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	go processJobs()

	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	log.Println("Fetching jobs")
	jobsList := make([]Job, len(jobs))
	for i, job := range jobs {
		jobsList[i] = *job
	}

	if len(jobsList) == 0 {
		log.Println("No jobs available")
	}
	err := json.NewEncoder(w).Encode(jobsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newJob Job
	err := json.NewDecoder(r.Body).Decode(&newJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// newJob.Duration = newJob.Duration * time.Second
	newJob.Status = "Pending"

	jobsMutex.Lock()
	heap.Push(&jobs, &newJob)
	jobsMutex.Unlock()

	// Wake up the processing goroutine if it's idle
	processingM.Lock()
	if !processing {
		processingM.Unlock()
		go processJobs()
	} else {
		processingM.Unlock()
	}

	err = json.NewEncoder(w).Encode(newJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func processJobs() {
	for {
		jobsMutex.Lock()
		if jobs.Len() == 0 {
			jobsMutex.Unlock()
			processingM.Lock()
			processing = false
			processingM.Unlock()
			return
		}

		job := heap.Pop(&jobs).(*Job)
		jobsMutex.Unlock()

		job.Status = "In Progress"
		log.Printf("Processing job: %+v\n", job)
		time.Sleep(job.Duration)

		job.Status = "Completed"
		log.Printf("Completed job: %+v\n", job)
	}
}
