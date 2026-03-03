package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ProjectsScreen struct {
	BaseScreen
	projectList  *widget.List
	currentTasks *widget.List
	selectedID   string
}

func NewProjectsScreen(service *gtd.Service, window fyne.Window) fyne.CanvasObject {
	screen := &ProjectsScreen{
		BaseScreen: BaseScreen{
			service: service,
			window:  window,
		},
	}

	return screen.buildUI()
}

func (s *ProjectsScreen) buildUI() fyne.CanvasObject {
	// ������ �������� (����� �� ��������, ������ �� ��������)
	s.projectList = s.createProjectList()

	// ������ ���������� �������
	addProjectBtn := widget.NewButtonWithIcon("����� ������", theme.ContentAddIcon(), func() {
		s.showNewProjectDialog()
	})

	// ������� ������ � �������
	s.currentTasks = s.createTaskList([]*gtd.Task{}, func(id string) {
		s.showTaskActionsInProject(id)
	})

	// ������ �������� ��� ��������
	projectActions := container.NewHBox(
		widget.NewButtonWithIcon("", theme.ConfirmIcon(), func() {
			if s.selectedID != "" {
				s.service.MoveToCompleted(s.selectedID)
				s.refreshProjectList()
				s.updateTasksForSelected()
			}
		}),
		widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			if s.selectedID != "" {
				s.service.MoveToTrash(s.selectedID)
				s.refreshProjectList()
				s.selectedID = ""
				s.updateTasksForSelected()
			}
		}),
	)

	// ���������� ���������
	addSubtask := widget.NewButtonWithIcon("�������� ���������", theme.ContentAddIcon(), func() {
		if s.selectedID != "" {
			s.showAddSubtaskDialog(s.selectedID)
		}
	})

	projectContent := container.NewBorder(
		container.NewVBox(
			addProjectBtn,
			widget.NewSeparator(),
			s.projectList,
		),
		nil, nil, nil,
		container.NewBorder(
			container.NewHBox(
				widget.NewLabel("������ �������:"),
				layout.NewSpacer(),
				projectActions,
			),
			addSubtask,
			nil, nil,
			container.NewScroll(s.currentTasks),
		),
	)

	split := container.NewHSplit(
		container.NewBorder(
			widget.NewLabelWithStyle("�������", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			nil, nil, nil,
			projectContent,
		),
		container.NewBorder(
			widget.NewLabelWithStyle("������", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			nil, nil, nil,
			container.NewCenter(widget.NewLabel("�������� ������")),
		),
	)
	split.Offset = 0.3

	return split
}

func (s *ProjectsScreen) createProjectList() *widget.List {
	projects, _ := s.service.GetProjects()

	list := widget.NewList(
		func() int {
			return len(projects) + 1 // +1 ��� ���������
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.FolderIcon()),
				widget.NewLabel("Project"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id == 0 {
				// ���������
				label := obj.(*fyne.Container).Objects[1].(*widget.Label)
				label.SetText("--- ��� ������� ---")
			} else {
				task := projects[id-1]
				label := obj.(*fyne.Container).Objects[1].(*widget.Label)
				label.SetText(task.Title)
			}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		if id == 0 {
			s.selectedID = ""
		} else {
			projects, _ := s.service.GetProjects()
			s.selectedID = projects[id-1].ID
		}
		s.updateTasksForSelected()
	}

	return list
}

func (s *ProjectsScreen) updateTasksForSelected() {
	if s.selectedID == "" {
		// ���������� ��� ������� ��� �����
		s.currentTasks.Hide()
		return
	}

	project, _ := s.service.GetTask(s.selectedID)
	if project != nil {
		s.currentTasks.Show()
		s.tasks = project.Subtasks
		s.currentTasks.Refresh()
	}
}

func (s *ProjectsScreen) refreshProjectList() {
	s.projectList.Refresh()
}

func (s *ProjectsScreen) showNewProjectDialog() {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("�������� �������")

	dialog.ShowCustomConfirm("����� ������", "�������", "������",
		container.NewVBox(
			widget.NewLabel("������� �������� �������:"),
			entry,
		),
		func(confirm bool) {
			if confirm && entry.Text != "" {
				s.service.AddToInbox(entry.Text, "")
				s.refreshProjectList()
			}
		},
		s.window,
	)
}

func (s *ProjectsScreen) showAddSubtaskDialog(projectID string) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("�������� ���������")

	dialog.ShowCustomConfirm("����� ���������", "��������", "������",
		container.NewVBox(
			widget.NewLabel("������� �������� ���������:"),
			entry,
		),
		func(confirm bool) {
			if confirm && entry.Text != "" {
				s.service.AddSubtask(projectID, entry.Text)
				s.updateTasksForSelected()
			}
		},
		s.window,
	)
}

func (s *ProjectsScreen) showTaskActionsInProject(taskID string) {
	s.showTaskActions(taskID)
}
