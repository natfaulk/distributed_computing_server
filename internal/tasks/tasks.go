package tasks

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/natfaulk/distributed_computing_server/internal/filehandling"
	"github.com/natfaulk/distributed_computing_server/internal/nflogger"
)

var logger *log.Logger = nflogger.Make("Tasks")

const (
	available string = "AVAILABLE"
	running   string = "RUNNING"
	complete  string = "COMPLETE"

	noTime    int64 = -1
	staleTime int64 = 2 * 60 * 60 // Two hours, should pull out to a config file
)

// Task holds all the info for a task
type Task struct {
	Params      string `json:"params"`
	ID          string `json:"id"`
	Status      string `json:"status"`
	Client      string `json:"client"`
	TimeAdded   int64  `json:"timeAdded"`
	TimeStarted int64  `json:"timeStarted"`
}

// Taskpool holds all Tasks
type Taskpool struct {
	Tasks []Task `json:"tasks"`
}

// NewTask creates and returns a new Task
func NewTask(_params string) Task {
	t := Task{_params, uuid.NewString(), available, "None", time.Now().Unix(), noTime}
	return t
}

// NewTaskpool creates a Taskpool and loads taskpool from file if it exists
func NewTaskpool() Taskpool {
	tp := Taskpool{}
	tp.LoadFromFile()
	return tp
}

// IsStale checks if job should have finished by now
func (task *Task) IsStale() bool {
	stale := (time.Now().Unix() - task.TimeStarted) > staleTime

	return task.Status == running && stale
}

// AddTask adds a task to the Taskpool
func (tp *Taskpool) AddTask(_params string) {
	t := NewTask(_params)
	tp.Tasks = append(tp.Tasks, t)
	tp.SaveToFile()
}

// GetTask gets and returns an available task if available,
// as well as a boolean indicating whether a task has been returned
func (tp *Taskpool) GetTask(_clientID string) (Task, bool) {
	for i, task := range tp.Tasks {
		if task.Status == available || task.IsStale() {
			tp.Tasks[i].Status = running
			tp.Tasks[i].Client = _clientID
			tp.Tasks[i].TimeStarted = time.Now().Unix()

			return tp.Tasks[i], true
		}
	}

	return Task{}, false
}

// GetTaskByID returns task by ID, returns the task (if found) and whether the task was found
func (tp *Taskpool) GetTaskByID(_taskID string) (Task, bool) {
	for _, task := range tp.Tasks {
		if task.ID == _taskID {
			return task, true
		}
	}

	return Task{}, false
}

// TaskComplete marks a task as completed
func (tp *Taskpool) TaskComplete(_taskID string, _clientID string) {
	for i, task := range tp.Tasks {
		if task.ID == _taskID {
			tp.Tasks[i].Status = complete
			tp.Tasks[i].Client = _clientID
			tp.SaveToFile()
			return
		}
	}
}

// SaveToFile saves current jobs to json file
func (tp Taskpool) SaveToFile() bool {
	tasksFilepath := filepath.Join(filehandling.DataFolder, filehandling.SaveFile)
	logger.Printf("Saving tasks file to %s", tasksFilepath)

	f, err := os.OpenFile(tasksFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		logger.Println("Error opening tasks file...", err)
		return false
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(tp)
	if err != nil {
		logger.Println("Error encoding tasks json file", err)
		return false
	}

	logger.Println("Successfully saved tasks json file")
	return true
}

// LoadFromFile loads a previously saved tasks file. If one does not exist, it creates an empty one
func (tp *Taskpool) LoadFromFile() bool {
	tasksFilepath := filepath.Join(filehandling.DataFolder, filehandling.SaveFile)
	logger.Printf("Loading tasks file from %s", tasksFilepath)

	f, err := os.Open(tasksFilepath)
	if os.IsNotExist(err) {
		logger.Println("Tasks file does not exist. Creating...")
		return tp.SaveToFile()
	} else if err != nil {
		logger.Println("Error opening tasks file", err)
		return false
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&tp)
	if err != nil {
		logger.Println("Error parsing tasks json file", err)
		return false
	}

	logger.Println("Successfully loaded tasks json file")

	return true
}
