package gtd

import (
	"errors"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddToInbox(title, voiceNotePath string) (*Task, error) {
	task, err := s.repo.CreateTask(title, Inbox)
	if err != nil {
		return nil, err
	}

	if voiceNotePath != "" {
		task.VoiceNote = voiceNotePath
		s.repo.UpdateTask(task)
	}

	return task, nil
}

func (s *Service) GetInbox() ([]*Task, error) {
	return s.repo.GetTasksByType(Inbox), nil
}

func (s *Service) GetProjects() ([]*Task, error) {
	return s.repo.GetTasksByType(Project), nil
}

func (s *Service) MoveToProject(taskID, projectID string) error {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return err
	}

	if task.Type == Inbox && projectID == "" {
		task.Type = Project
		return s.repo.UpdateTask(task)
	}

	if projectID != "" {
		parent, err := s.repo.GetTask(projectID)
		if err != nil {
			return err
		}

		task.ParentID = projectID
		parent.Subtasks = append(parent.Subtasks, task)
		s.repo.UpdateTask(parent)
	}

	return s.repo.UpdateTask(task)
}

func (s *Service) MoveToCompleted(taskID string) error {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return err
	}

	now := time.Now()
	task.Status = Completed
	task.CompletedAt = &now

	return s.repo.UpdateTask(task)
}

func (s *Service) MoveToTrash(taskID string) error {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return err
	}

	task.Status = Trash
	return s.repo.UpdateTask(task)
}

func (s *Service) AddSubtask(parentID, title string) (*Task, error) {
	parent, err := s.repo.GetTask(parentID)
	if err != nil {
		return nil, err
	}

	if parent.Type != Project {
		return nil, errors.New("only projects can have subtasks")
	}

	subtask, err := s.repo.CreateTask(title, Project)
	if err != nil {
		return nil, err
	}

	subtask.ParentID = parentID
	parent.Subtasks = append(parent.Subtasks, subtask)

	if err := s.repo.UpdateTask(parent); err != nil {
		return nil, err
	}

	return subtask, nil
}

func (s *Service) DeletePermanently(taskID string) error {
	return s.repo.DeleteTask(taskID)
}

func (s *Service) RestoreFromTrash(taskID string) error {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return err
	}

	task.Status = ""
	return s.repo.UpdateTask(task)
}

func (s *Service) UpdateTask(task *Task) error {
	return s.repo.UpdateTask(task)
}

func (s *Service) GetTask(id string) (*Task, error) {
	return s.repo.GetTask(id)
}

func (s *Service) GetTasksByStatus(status TaskStatus) []*Task {
	return s.repo.GetTasksByStatus(status)
}

func (s *Service) ProcessVoiceInput(audioData []byte) (string, error) {
	// Заглушка для голосового ввода
	return "Распознанный текст из голосового сообщения", nil
}
