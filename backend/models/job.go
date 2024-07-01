package models

import "time"

// Job represents data about a posted job
type Job struct {
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
	Status   string        `json:"status"`
	Index    int           `json:"-"`  // Index in the heap
	ID       int           `json:"id"` // Unique identifier for the job
}
