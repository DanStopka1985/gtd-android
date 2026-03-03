package gtd

import (
	"time"
)

type TaskType string
type TaskStatus string

const (
	Inbox     TaskType   = "inbox"
	Project   TaskType   = "project"
	Completed TaskStatus = "completed"
	Trash     TaskStatus = "trash"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        TaskType   `json:"type"`
	Status      TaskStatus `json:"status"`
	ParentID    string     `json:"parent_id,omitempty"`
	Subtasks    []*Task    `json:"subtasks,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	VoiceNote   string     `json:"voice_note,omitempty"`
}

type GTDService interface {
	AddToInbox(title, voiceNotePath string) (*Task, error)
	GetInbox() ([]*Task, error)
	GetProjects() ([]*Task, error)
	GetTasksByStatus(status TaskStatus) []*Task
	MoveToProject(taskID, projectID string) error
	MoveToCompleted(taskID string) error
	MoveToTrash(taskID string) error
	AddSubtask(parentID, title string) (*Task, error)
	DeletePermanently(taskID string) error
	RestoreFromTrash(taskID string) error
	UpdateTask(task *Task) error
	GetTask(id string) (*Task, error)
	ProcessVoiceInput(audioData []byte) (string, error)
}
