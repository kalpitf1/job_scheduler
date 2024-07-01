package utils

import (
	"container/heap"
	"sync"

	"github.com/kalpitf1/job_scheduler/backend/models"
)

// JobPriorityQueue implements heap.Interface and holds Jobs
type JobPriorityQueue []*models.Job

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
	job := x.(*models.Job)
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

func (pq *JobPriorityQueue) Peek() *models.Job {
	if len(*pq) == 0 {
		return nil
	}
	return (*pq)[0]
}

var (
	pq      JobPriorityQueue
	pqMutex sync.Mutex
)

func init() {
	heap.Init(&pq)
}

func PushJob(job *models.Job) {
	pqMutex.Lock()
	heap.Push(&pq, job)
	pqMutex.Unlock()
}

func PopJob() *models.Job {
	pqMutex.Lock()
	defer pqMutex.Unlock()
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(&pq).(*models.Job)
}
