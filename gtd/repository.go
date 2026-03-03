package gtd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	mu      sync.RWMutex
	dataDir string
	tasks   map[string]*Task
}

func NewRepository(dataDir string) (*Repository, error) {
	repo := &Repository{
		dataDir: dataDir,
		tasks:   make(map[string]*Task),
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) load() error {
	file := filepath.Join(r.dataDir, "tasks.json")
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var tasks []*Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, task := range tasks {
		r.tasks[task.ID] = task
	}

	return nil
}

func (r *Repository) save() error {
	r.mu.RLock()
	tasks := make([]*Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	r.mu.RUnlock()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	file := filepath.Join(r.dataDir, "tasks.json")
	return os.WriteFile(file, data, 0644)
}

func (r *Repository) CreateTask(title string, taskType TaskType) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task := &Task{
		ID:        uuid.New().String(),
		Title:     title,
		Type:      taskType,
		Status:    "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Subtasks:  []*Task{},
	}

	r.tasks[task.ID] = task
	return task, r.save()
}

func (r *Repository) GetTask(id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, os.ErrNotExist
	}
	return task, nil
}

func (r *Repository) UpdateTask(task *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task
	return r.save()
}

func (r *Repository) DeleteTask(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tasks, id)
	return r.save()
}

func (r *Repository) GetTasksByType(taskType TaskType) []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Task
	for _, task := range r.tasks {
		if task.Type == taskType && task.Status != Trash {
			result = append(result, task)
		}
	}
	return result
}

func (r *Repository) GetTasksByStatus(status TaskStatus) []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Task
	for _, task := range r.tasks {
		if task.Status == status {
			result = append(result, task)
		}
	}
	return result
}
