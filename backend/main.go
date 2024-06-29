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
	"github.com/gorilla/websocket"
)

// Job represents data about a posted job
type Job struct {
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
	Status   string        `json:"status"`
	Index    int           `json:"-"`  // Index in the heap
	ID       int           `json:"id"` // Unique identifier for the job
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
	jobs        []*Job
	jobsMutex   sync.Mutex
	jobCounter  int
	pq          JobPriorityQueue
	pqMutex     sync.Mutex
	processing  bool
	processingM sync.Mutex

	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan *Job)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	log.Println("Starting Backend")

	heap.Init(&pq)

	r := mux.NewRouter()
	r.HandleFunc("/jobs", getJobs).Methods("GET")
	r.HandleFunc("/jobs", createJob).Methods("POST")
	r.HandleFunc("/ws", handleConnections)

	// Apply CORS headers to the router
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Allow requests from your React app
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	go handleMessages()

	log.Fatal(http.ListenAndServe(":8080", cors(r)))
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	err := json.NewEncoder(w).Encode(jobs)
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

	newJob.Status = "Pending"
	newJob.ID = jobCounter

	jobsMutex.Lock()
	jobCounter++
	jobs = append(jobs, &newJob)
	jobsMutex.Unlock()

	pqMutex.Lock()
	heap.Push(&pq, &newJob)
	pqMutex.Unlock()

	broadcast <- &newJob

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

		job := heap.Pop(&pq).(*Job)
		pqMutex.Unlock()

		jobsMutex.Lock()
		job.Status = "In Progress"
		jobsMutex.Unlock()
		log.Printf("Processing job: %+v\n", job)

		broadcast <- job

		time.Sleep(job.Duration)

		jobsMutex.Lock()
		job.Status = "Completed"
		jobsMutex.Unlock()
		log.Printf("Completed job: %+v\n", job)

		broadcast <- job
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Job
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		job := <-broadcast

		for client := range clients {
			err := client.WriteJSON(job)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
