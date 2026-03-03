package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"gtd-android/gtd"
)

type BaseScreen struct {
	service *gtd.Service
	window  fyne.Window
	list    *widget.List
	tasks   []*gtd.Task
}

func (s *BaseScreen) refreshList() {}

func (s *BaseScreen) showTaskActions(task *gtd.Task) {
	log.Println("Showing actions for task:", task.Title)

	content := container.NewVBox(
		widget.NewButton("✏️ Редактировать", func() {
			s.showEditDialog(task)
		}),
		widget.NewButton("📋 Декомпозировать", func() {
			s.showDecomposeDialog(task)
		}),
		widget.NewButton("🗑 Удалить", func() {
			s.service.MoveToTrash(task.ID)
			s.refreshList()
		}),
	)

	dialog.ShowCustom("Действия", "Закрыть", content, s.window)
}

func (s *BaseScreen) showEditDialog(task *gtd.Task) {
	entry := widget.NewEntry()
	entry.SetText(task.Title)

	content := container.NewVBox(
		widget.NewLabel("Новое название:"),
		entry,
		widget.NewButton("Сохранить", func() {
			if entry.Text != "" {
				task.Title = entry.Text
				s.service.UpdateTask(task)
				s.refreshList()
			}
		}),
	)

	dialog.ShowCustom("Редактировать", "Отмена", content, s.window)
}

func (s *BaseScreen) showDecomposeDialog(parent *gtd.Task) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Название подзадачи")

	content := container.NewVBox(
		widget.NewLabel("Новая подзадача:"),
		entry,
		widget.NewButton("Добавить", func() {
			if entry.Text != "" {
				s.service.AddSubtask(parent.ID, entry.Text)
				s.refreshList()
			}
		}),
	)

	dialog.ShowCustom("Декомпозировать", "Отмена", content, s.window)
}

func (s *BaseScreen) createTaskList(items []*gtd.Task, onSelect func(id string)) *widget.List {
	s.tasks = items

	list := widget.NewList(
		func() int {
			return len(s.tasks)
		},
		func() fyne.CanvasObject {
			// Максимально простой шаблон
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("Task"),
				widget.NewButton("⋮", nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(s.tasks) {
				task := s.tasks[id]
				box := obj.(*fyne.Container)

				// Обновляем label
				if label, ok := box.Objects[1].(*widget.Label); ok {
					label.SetText(task.Title)
				}

				// Обновляем кнопку
				if btn, ok := box.Objects[2].(*widget.Button); ok {
					btn.OnTapped = func() {
						s.showTaskActions(task)
					}
				}
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if onSelect != nil && id < len(s.tasks) {
			onSelect(s.tasks[id].ID)
		}
		list.UnselectAll()
	}

	return list
}
