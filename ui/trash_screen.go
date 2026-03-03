package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"gtd-android/gtd"
)

type TrashScreen struct {
	BaseScreen
}

func NewTrashScreen(service *gtd.Service, window fyne.Window) fyne.CanvasObject {
	screen := &TrashScreen{
		BaseScreen: BaseScreen{
			service: service,
			window:  window,
		},
	}

	return screen.buildUI()
}

func (s *TrashScreen) buildUI() fyne.CanvasObject {
	restoreAllBtn := widget.NewButtonWithIcon("Восстановить все", theme.ViewRestoreIcon(), func() {
		s.showRestoreAllDialog()
	})

	emptyBtn := widget.NewButtonWithIcon("Очистить корзину", theme.DeleteIcon(), func() {
		s.showEmptyTrashDialog()
	})

	s.list = s.createTrashList()

	return container.NewBorder(
		container.NewHBox(
			widget.NewLabel("Корзина"),
			restoreAllBtn,
			emptyBtn,
		),
		nil, nil, nil,
		container.NewScroll(s.list),
	)
}

func (s *TrashScreen) createTrashList() *widget.List {
	tasks := s.service.GetTasksByStatus(gtd.Trash)

	list := widget.NewList(
		func() int {
			return len(tasks)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DeleteIcon()),
				widget.NewLabel("Task"),
				widget.NewButtonWithIcon("", theme.ViewRestoreIcon(), nil),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < len(tasks) {
				task := tasks[id]
				box := obj.(*fyne.Container)
				label := box.Objects[1].(*widget.Label)
				label.SetText(task.Title)

				restoreBtn := box.Objects[2].(*widget.Button)
				restoreBtn.OnTapped = func() {
					s.service.RestoreFromTrash(task.ID)
					s.refreshList()
				}

				deleteBtn := box.Objects[3].(*widget.Button)
				deleteBtn.OnTapped = func() {
					s.showDeleteConfirmDialog(task.ID)
				}
			}
		},
	)

	return list
}

func (s *TrashScreen) refreshList() {
	s.list = s.createTrashList()
	s.list.Refresh()
}

func (s *TrashScreen) showRestoreAllDialog() {
	dialog.ShowConfirm("Восстановить все",
		"Восстановить все задачи из корзины?",
		func(confirm bool) {
			if confirm {
				tasks := s.service.GetTasksByStatus(gtd.Trash)
				for _, task := range tasks {
					s.service.RestoreFromTrash(task.ID)
				}
				s.refreshList()
			}
		},
		s.window,
	)
}

func (s *TrashScreen) showEmptyTrashDialog() {
	dialog.ShowConfirm("Очистить корзину",
		"Все задачи будут удалены безвозвратно. Продолжить?",
		func(confirm bool) {
			if confirm {
				tasks := s.service.GetTasksByStatus(gtd.Trash)
				for _, task := range tasks {
					s.service.DeletePermanently(task.ID)
				}
				s.refreshList()
			}
		},
		s.window,
	)
}

func (s *TrashScreen) showDeleteConfirmDialog(taskID string) {
	dialog.ShowConfirm("Удалить навсегда",
		"Задача будет удалена безвозвратно. Продолжить?",
		func(confirm bool) {
			if confirm {
				s.service.DeletePermanently(taskID)
				s.refreshList()
			}
		},
		s.window,
	)
}
