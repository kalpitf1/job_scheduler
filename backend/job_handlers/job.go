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
	Jobs                  []*models.Job
	JobsMutex             sync.Mutex
	jobCounter            int
	pq                    utils.JobPriorityQueue
	pqMutex               sync.Mutex
	processing            bool
	processingM           sync.Mutex
	BroadcastEnabled      = true
	BroadcastEnabledMutex sync.Mutex
)

func init() {
	pq = make(utils.JobPriorityQueue, 0)
	heap.Init(&pq)
}

func Reset() {
	pqMutex.Lock()
	defer pqMutex.Unlock()
	pq = make(utils.JobPriorityQueue, 0)
	heap.Init(&pq)
	Jobs = make([]*models.Job, 0)
	jobCounter = 0
	processing = false
	BroadcastEnabled = false
}

func GetJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	JobsMutex.Lock()
	defer JobsMutex.Unlock()

	err := json.NewEncoder(w).Encode(Jobs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Println("CreateJob: Start")

	var newJob models.Job
	err := json.NewDecoder(r.Body).Decode(&newJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		// log.Println("CreateJob: Error decoding request body:", err)
		return
	}

	// log.Printf("CreateJob: Received job - %+v\n", newJob)

	newJob.Status = "Pending"
	newJob.ID = jobCounter

	JobsMutex.Lock()
	jobCounter++
	Jobs = append(Jobs, &newJob)
	JobsMutex.Unlock()

	pqMutex.Lock()
	heap.Push(&pq, &newJob)
	pqMutex.Unlock()

	BroadcastEnabledMutex.Lock()
	// log.Printf("CreateJob: BroadcastEnabled is %v\n", BroadcastEnabled)
	// broadcast enabled by default, can be turned off for testing
	if BroadcastEnabled {
		// websocket.Broadcast <- &newJob

		select {
		case websocket.Broadcast <- &newJob:
			log.Println("CreateJob: Broadcast successful")
		default:
			log.Println("CreateJob: Broadcast channel is full, skipping broadcast")
		}
	}
	BroadcastEnabledMutex.Unlock()

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
		// log.Println("CreateJob: Error encoding response:", err)
		return
	}

	// log.Println("CreateJob: End")
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

		// BroadcastEnabledMutex.Lock()
		// log.Printf("BroadcastEnabled before heap pop: %v", BroadcastEnabled)
		// BroadcastEnabledMutex.Unlock()
		job := heap.Pop(&pq).(*models.Job)
		// BroadcastEnabledMutex.Lock()
		// log.Printf("BroadcastEnabled after heap pop: %v", BroadcastEnabled)
		// BroadcastEnabledMutex.Unlock()
		pqMutex.Unlock()

		JobsMutex.Lock()
		// BroadcastEnabledMutex.Lock()
		// log.Printf("BroadcastEnabled before changing job status to In Progress: %v", BroadcastEnabled)
		// BroadcastEnabledMutex.Unlock()
		job.Status = "In Progress"
		JobsMutex.Unlock()
		log.Printf("Processing job: %+v\n", job)

		BroadcastEnabledMutex.Lock()
		// broadcast enabled by default, can be turned off for testing
		if BroadcastEnabled {
			// websocket.Broadcast <- job

			select {
			case websocket.Broadcast <- job:
				log.Println("CreateJob: Broadcast successful")
			default:
				log.Println("CreateJob: Broadcast channel is full, skipping broadcast")
			}
		}
		BroadcastEnabledMutex.Unlock()

		// log.Printf("Sleeping for %+v", job.Duration)
		time.Sleep(job.Duration)
		// log.Printf("Done sleeping for %+v", job.Duration)

		JobsMutex.Lock()
		// BroadcastEnabledMutex.Lock()
		// log.Printf("BroadcastEnabled before changing job status to Completed: %v", BroadcastEnabled)
		// BroadcastEnabledMutex.Unlock()
		job.Status = "Completed"
		JobsMutex.Unlock()
		log.Printf("Completed job: %+v\n", job)

		BroadcastEnabledMutex.Lock()
		// broadcast enabled by default, can be turned off for testing
		if BroadcastEnabled {
			// websocket.Broadcast <- job

			select {
			case websocket.Broadcast <- job:
				log.Println("CreateJob: Broadcast successful")
			default:
				log.Println("CreateJob: Broadcast channel is full, skipping broadcast")
			}
		}
		BroadcastEnabledMutex.Unlock()
	}
}
