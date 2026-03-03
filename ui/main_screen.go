package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type BaseScreen struct {
	service *gtd.Service
	window  fyne.Window
	list    *widget.List
	tasks   []*gtd.Task
}

func (s *BaseScreen) refreshList() {
	// ����� ��������������
}

func (s *BaseScreen) showTaskActions(task *gtd.Task) {
	popup := widget.NewPopUp(
		container.NewVBox(
			widget.NewButton("�������������", func() {
				s.showEditDialog(task)
			}),
			widget.NewButton("���������������", func() {
				s.showDecomposeDialog(task)
			}),
			widget.NewButton("�������", func() {
				s.service.MoveToTrash(task.ID)
				s.refreshList()
			}),
			widget.NewButton("������", func() {
				popup.Hide()
			}),
		),
		s.window.Canvas(),
	)

	popup.ShowAtPosition(fyne.NewPos(100, 200))
}

func (s *BaseScreen) showEditDialog(task *gtd.Task) {
	entry := widget.NewMultiLineEntry()
	entry.SetText(task.Title)

	dialog.ShowCustomConfirm("������������� ������", "���������", "������",
		container.NewVBox(
			widget.NewLabel("��������:"),
			entry,
		),
		func(confirm bool) {
			if confirm && entry.Text != "" {
				task.Title = entry.Text
				s.service.UpdateTask(task)
				s.refreshList()
			}
		},
		s.window,
	)
}

func (s *BaseScreen) showDecomposeDialog(parent *gtd.Task) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("�������� ���������")

	dialog.ShowCustomConfirm("��������������� ������", "��������", "������",
		container.NewVBox(
			widget.NewLabel("����� ���������:"),
			entry,
		),
		func(confirm bool) {
			if confirm && entry.Text != "" {
				s.service.AddSubtask(parent.ID, entry.Text)
				s.refreshList()
			}
		},
		s.window,
	)
}

func (s *BaseScreen) createTaskList(items []*gtd.Task, onSelect func(id string)) *widget.List {
	s.tasks = items

	list := widget.NewList(
		func() int {
			return len(s.tasks)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("Task"),
				widget.NewButtonWithIcon("", theme.MoreHorizIcon(), func() {}),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			task := s.tasks[id]
			box := obj.(*fyne.Container)
			label := box.Objects[1].(*widget.Label)
			label.SetText(task.Title)

			// ������ ��������
			btn := box.Objects[2].(*widget.Button)
			btn.OnTapped = func() {
				s.showTaskActions(task)
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if onSelect != nil {
			onSelect(s.tasks[id].ID)
		}
		list.UnselectAll()
	}

	return list
}
