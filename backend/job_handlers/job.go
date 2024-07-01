package job_handlers

import (
	"container/heap"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kalpitf1/job_scheduler/backend/models"
	"github.com/kalpitf1/job_scheduler/backend/utils"
	"github.com/kalpitf1/job_scheduler/backend/websocket"
)

var (
	jobs        []*models.Job
	jobsMutex   sync.Mutex
	jobCounter  int
	pq          utils.JobPriorityQueue
	pqMutex     sync.Mutex
	processing  bool
	processingM sync.Mutex
)

func init() {
	heap.Init(&pq)
}

func GetJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	err := json.NewEncoder(w).Encode(jobs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newJob models.Job
	err := json.NewDecoder(r.Body).Decode(&newJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newJob.Status = "Pending"
	newJob.ID = jobCounter

	jobsMutex.Lock()
	jobCounter++
	jobs = append(jobs, &newJob)
	jobsMutex.Unlock()

	pqMutex.Lock()
	heap.Push(&pq, &newJob)
	pqMutex.Unlock()

	websocket.Broadcast <- &newJob

	// Wake up the processing goroutine if it's idle
	processingM.Lock()
	if !processing {
		processing = true
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
		pqMutex.Lock()
		if pq.Len() == 0 {
			pqMutex.Unlock()
			processingM.Lock()
			processing = false
			processingM.Unlock()
			return
		}

		job := heap.Pop(&pq).(*models.Job)
		pqMutex.Unlock()

		jobsMutex.Lock()
		job.Status = "In Progress"
		jobsMutex.Unlock()
		log.Printf("Processing job: %+v\n", job)

		websocket.Broadcast <- job

		time.Sleep(job.Duration)

		jobsMutex.Lock()
		job.Status = "Completed"
		jobsMutex.Unlock()
		log.Printf("Completed job: %+v\n", job)

		websocket.Broadcast <- job
	}
}
