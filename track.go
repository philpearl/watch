package rebuilder

import (
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/philpearl/rebuilder/wire"
)

type taskKey struct {
	name     string
	taskType wire.TaskType
}

type Track struct {
	sync.Mutex
	tasks map[taskKey]*wire.Task
	list  []*wire.Task
}

func NewTrack() *Track {
	return &Track{
		tasks: make(map[taskKey]*wire.Task),
	}
}

func (t *Track) get(name string, taskType wire.TaskType) *wire.Task {
	key := taskKey{name: name, taskType: taskType}
	status, ok := t.tasks[key]
	if !ok {
		status = &wire.Task{
			Name: name,
			Type: taskType,
		}
		t.tasks[key] = status
		t.list = append(t.list, status)
	}

	return status
}

func (t *Track) Pending(name string, taskType wire.TaskType) {
	t.Lock()
	defer t.Unlock()
	status := t.get(name, taskType)
	status.Status = wire.StatusPending
}

func (t *Track) Started(name string, taskType wire.TaskType) {
	t.Lock()
	defer t.Unlock()

	status := t.get(name, taskType)
	status.Status = wire.StatusRunning
	status.Started = time.Now()
	status.Ended = time.Time{}
}

func (t *Track) Ended(name string, taskType wire.TaskType, err error, output string) {
	t.Lock()
	defer t.Unlock()

	status := t.get(name, taskType)
	status.Status = wire.StatusComplete
	status.Ended = time.Now()
	status.Error = err
	status.Output = output
}

func (t *Track) WriteStatus(w io.Writer) error {
	t.Lock()
	defer t.Unlock()
	return json.NewEncoder(w).Encode(t.list)
}
