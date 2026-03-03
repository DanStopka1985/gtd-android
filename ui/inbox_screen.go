package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type InboxScreen struct {
	BaseScreen
	voiceProc *gtd.VoiceProcessor
	addBtn    *widget.Button
	voiceBtn  *widget.Button
}

func NewInboxScreen(service *gtd.Service, voiceProc *gtd.VoiceProcessor, window fyne.Window) fyne.CanvasObject {
	screen := &InboxScreen{
		BaseScreen: BaseScreen{
			service: service,
			window:  window,
		},
		voiceProc: voiceProc,
	}

	return screen.buildUI()
}

func (s *InboxScreen) buildUI() fyne.CanvasObject {
	// ���� ����� ��� �������� ����������
	input := widget.NewEntry()
	input.SetPlaceHolder("����� ������...")

	s.addBtn = widget.NewButtonWithIcon("��������", theme.ContentAddIcon(), func() {
		if input.Text != "" {
			s.service.AddToInbox(input.Text, "")
			input.SetText("")
			s.refreshList()
		}
	})

	// ������ ���������� �����
	s.voiceBtn = widget.NewButtonWithIcon("", theme.MediaRecordIcon(), func() {
		s.showVoiceInputDialog()
	})

	input.OnSubmitted = func(text string) {
		if text != "" {
			s.service.AddToInbox(text, "")
			input.SetText("")
			s.refreshList()
		}
	}

	// ������ �����
	s.list = s.createTaskList(s.getTasks(), func(id string) {
		s.showMoveToProjectDialog(id)
	})

	header := container.NewBorder(
		nil, nil, nil,
		container.NewHBox(s.voiceBtn, s.addBtn),
		input,
	)

	return container.NewBorder(
		header,
		nil, nil, nil,
		container.NewScroll(s.list),
	)
}

func (s *InboxScreen) getTasks() []*gtd.Task {
	tasks, _ := s.service.GetInbox()
	return tasks
}

func (s *InboxScreen) refreshList() {
	s.tasks = s.getTasks()
	s.list.Refresh()
}

func (s *InboxScreen) showVoiceInputDialog() {
	dialog.ShowInformation("��������� ����",
		"������� '������ ������' � ����������� ������",
		s.window)

	// � �������� ���������� ����� ����� ����� Android SpeechRecognizer
	// ����� gobind ��� �������� ������

	go func() {
		// �������� ������������� ������
		// � ���������� ����� ����� ���������� � Android SpeechRecognizer
		recognizedText := "������ �������� �� ������"

		s.window.Canvas().Overlays().RemoveAll()
		s.service.AddToInbox(recognizedText, "")
		s.refreshList()
	}()
}

func (s *InboxScreen) showMoveToProjectDialog(taskID string) {
	projects, _ := s.service.GetProjects()

	var items []string
	items = append(items, "[����� ������]")
	for _, p := range projects {
		items = append(items, p.Title)
	}

	dialog.ShowCustom("����������� � ������", "������",
		container.NewVBox(
			widget.NewLabel("�������� ������:"),
			widget.NewSelect(items, func(selected string) {
				if selected == "[����� ������]" {
					s.showNewProjectDialog(taskID)
				} else {
					for _, p := range projects {
						if p.Title == selected {
							s.service.MoveToProject(taskID, p.ID)
							s.refreshList()
							break
						}
					}
				}
			}),
		),
		s.window,
	)
}

func (s *InboxScreen) showNewProjectDialog(taskID string) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("�������� �������")

	dialog.ShowCustomConfirm("����� ������", "�������", "������",
		container.NewVBox(
			widget.NewLabel("������� �������� �������:"),
			entry,
		),
		func(confirm bool) {
			if confirm && entry.Text != "" {
				// ������� ����� ������ � ���������� ������
				project, _ := s.service.AddToInbox(entry.Text, "")
				s.service.MoveToProject(project.ID, "")
				s.service.MoveToProject(taskID, project.ID)
				s.refreshList()
			}
		},
		s.window,
	)
}
