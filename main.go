package main

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/DanStopka/gtd-android/gtd"
	"github.com/DanStopka/gtd-android/ui"
)

type GTDApp struct {
	service    *gtd.Service
	repo       *gtd.Repository
	voiceProc  *gtd.VoiceProcessor
	mainWindow fyne.Window
}

func main() {
	// ������������� ���������� Fyne
	myApp := app.NewWithID("com.gtdandroid.app")
	myApp.Settings().SetTheme(theme.DarkTheme())

	window := myApp.NewWindow("GTD Organizer")
	window.Resize(fyne.NewSize(360, 640))

	// ������������� ���������
	dataDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	dataDir = filepath.Join(dataDir, "gtd-android")

	repo, err := gtd.NewRepository(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	service := gtd.NewService(repo)
	voiceProc := gtd.NewVoiceProcessor(dataDir)

	app := &GTDApp{
		service:    service,
		repo:       repo,
		voiceProc:  voiceProc,
		mainWindow: window,
	}

	app.setupUI()
	window.ShowAndRun()
}

func (a *GTDApp) setupUI() {
	// ������� ������� ��� ���������
	inboxTab := container.NewTabItemWithIcon("��������", theme.InboxIcon(), ui.NewInboxScreen(a.service, a.voiceProc, a.mainWindow))
	projectsTab := container.NewTabItemWithIcon("�������", theme.FolderIcon(), ui.NewProjectsScreen(a.service, a.mainWindow))
	completedTab := container.NewTabItemWithIcon("�����������", theme.ConfirmIcon(), ui.NewCompletedScreen(a.service, a.mainWindow))
	trashTab := container.NewTabItemWithIcon("�������", theme.DeleteIcon(), ui.NewTrashScreen(a.service, a.mainWindow))

	tabs := container.NewAppTabs(
		inboxTab,
		projectsTab,
		completedTab,
		trashTab,
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	a.mainWindow.SetContent(tabs)
}
