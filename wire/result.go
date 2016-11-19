package wire

import (
	"time"
)

type TaskType string

const (
	TaskTypeTest  TaskType = "test"
	TaskTypeBuild          = "build"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusRunning         = "running"
	StatusComplete        = "complete"
)

type Task struct {
	Name    string    `json:"name"`
	Type    TaskType  `json:"type"`
	Status  Status    `json:"status"`
	Started time.Time `json:"started,omitempty"`
	Ended   time.Time `json:"ended,omitempty"`
	Output  string    `json:"output"`
	Error   error     `json:"error"`
}
