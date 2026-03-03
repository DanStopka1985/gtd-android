package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"gtd-android/gtd"
)

type CompletedScreen struct {
	BaseScreen
}

func NewCompletedScreen(service *gtd.Service, window fyne.Window) fyne.CanvasObject {
	screen := &CompletedScreen{
		BaseScreen: BaseScreen{
			service: service,
			window:  window,
		},
	}

	return screen.buildUI()
}

func (s *CompletedScreen) buildUI() fyne.CanvasObject {
	clearBtn := widget.NewButtonWithIcon("Очистить все", theme.DeleteIcon(), func() {
		tasks := s.getTasks()
		for _, task := range tasks {
			s.service.DeletePermanently(task.ID)
		}
		s.refreshList()
	})

	s.list = s.createTaskList(s.getTasks(), nil)

	return container.NewBorder(
		container.NewHBox(
			widget.NewLabel("Выполненные задачи"),
			clearBtn,
		),
		nil, nil, nil,
		container.NewScroll(s.list),
	)
}

func (s *CompletedScreen) getTasks() []*gtd.Task {
	return s.service.GetTasksByStatus(gtd.Completed)
}

func (s *CompletedScreen) refreshList() {
	s.tasks = s.getTasks()
	if s.list != nil {
		s.list.Refresh()
	}
}
