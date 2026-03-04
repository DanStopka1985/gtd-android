package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Task struct {
	Title    string
	Done     bool
	Subtasks []Task
}

type Data struct {
	Inbox     []Task
	Projects  []Task
	Completed []Task
	Trash     []Task
}

func main() {
	// Для отладки
	f, err := os.Create("/sdcard/gtd_debug.log")
	if err == nil {
		log.SetOutput(f)
	}
	log.Println("=== GTD App Started ===")

	myApp := app.New()
	myWindow := myApp.NewWindow("GTD Organizer")
	myWindow.Resize(fyne.NewSize(500, 600))

	// Загружаем данные
	dataDir, err := os.UserCacheDir()
	if err != nil {
		log.Println("Error getting cache dir:", err)
		dataDir = os.TempDir()
	}
	dataDir = filepath.Join(dataDir, "gtd-data")
	os.MkdirAll(dataDir, 0755)
	dataFile := filepath.Join(dataDir, "tasks.json")
	log.Println("Data file:", dataFile)

	data := &Data{}

	// Читаем существующие данные
	if file, err := os.ReadFile(dataFile); err == nil {
		err = json.Unmarshal(file, data)
		if err != nil {
			log.Println("Error parsing JSON:", err)
		} else {
			log.Printf("Loaded: Inbox=%d, Projects=%d, Completed=%d, Trash=%d",
				len(data.Inbox), len(data.Projects), len(data.Completed), len(data.Trash))
		}
	} else {
		log.Println("No existing data, creating test data")
		data.Inbox = []Task{
			{Title: "Купить продукты", Done: false},
			{Title: "Сделать зарядку", Done: false},
			{Title: "Позвонить маме", Done: false},
		}
		data.Projects = []Task{
			{
				Title: "Учеба",
				Done:  false,
				Subtasks: []Task{
					{Title: "Прочитать главу 5", Done: false},
					{Title: "Сделать конспект", Done: false},
				},
			},
		}
		// Сразу сохраняем тестовые данные
		saveDataToFile(dataFile, data)
	}

	currentView := "inbox"

	// Заголовок
	title := widget.NewLabelWithStyle("GTD Organizer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Счетчик задач
	counter := widget.NewLabel("")
	updateCounter := func() {
		switch currentView {
		case "inbox":
			counter.SetText(fmt.Sprintf("📥 %d задач", len(data.Inbox)))
		case "projects":
			counter.SetText(fmt.Sprintf("📁 %d проектов", len(data.Projects)))
		case "completed":
			counter.SetText(fmt.Sprintf("✅ %d выполнено", len(data.Completed)))
		case "trash":
			counter.SetText(fmt.Sprintf("🗑 %d в корзине", len(data.Trash)))
		}
	}

	// Список задач
	taskList := widget.NewList(
		func() int {
			switch currentView {
			case "inbox":
				return len(data.Inbox)
			case "projects":
				return len(data.Projects)
			case "completed":
				return len(data.Completed)
			case "trash":
				return len(data.Trash)
			default:
				return 0
			}
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(nil),
				widget.NewLabel(""),
			)
		},
		func(i int, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			icon := box.Objects[0].(*widget.Icon)
			label := box.Objects[1].(*widget.Label)

			switch currentView {
			case "inbox":
				icon.SetResource(theme.MailComposeIcon())
				label.SetText(data.Inbox[i].Title)
			case "projects":
				icon.SetResource(theme.FolderIcon())
				label.SetText(data.Projects[i].Title)
			case "completed":
				icon.SetResource(theme.ConfirmIcon())
				label.SetText(data.Completed[i].Title)
			case "trash":
				icon.SetResource(theme.DeleteIcon())
				label.SetText(data.Trash[i].Title)
			}
		},
	)

	// Общая функция обновления
	refreshAll := func() {
		saveDataToFile(dataFile, data)
		taskList.Refresh()
		updateCounter()
	}

	taskList.OnSelected = func(id widget.ListItemID) {
		switch currentView {
		case "inbox":
			showInboxActions(&data.Inbox[id], data, id, myWindow, refreshAll)
		case "projects":
			showProjectActions(&data.Projects[id], data, id, myWindow, refreshAll)
		case "completed":
			showCompletedActions(&data.Completed[id], data, id, myWindow, refreshAll)
		case "trash":
			showTrashActions(&data.Trash[id], data, id, myWindow, refreshAll)
		}
		taskList.UnselectAll()
	}

	// Поле ввода
	input := widget.NewEntry()
	input.SetPlaceHolder("➕ Новая задача...")
	input.OnSubmitted = func(text string) {
		if text != "" {
			switch currentView {
			case "inbox":
				data.Inbox = append(data.Inbox, Task{Title: text, Done: false})
				refreshAll()
				input.SetText("")
			case "projects":
				data.Projects = append(data.Projects, Task{Title: text, Done: false, Subtasks: []Task{}})
				refreshAll()
				input.SetText("")
			}
		}
	}

	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if input.Text != "" {
			switch currentView {
			case "inbox":
				data.Inbox = append(data.Inbox, Task{Title: input.Text, Done: false})
				refreshAll()
				input.SetText("")
			case "projects":
				data.Projects = append(data.Projects, Task{Title: input.Text, Done: false, Subtasks: []Task{}})
				refreshAll()
				input.SetText("")
			}
		}
	})

	inputRow := container.NewBorder(nil, nil, nil, addBtn, input)

	// Кнопки навигации
	inboxBtn := widget.NewButtonWithIcon("Входящие", theme.MailComposeIcon(), func() {
		currentView = "inbox"
		taskList.Refresh()
		input.Show()
		addBtn.Show()
		input.SetPlaceHolder("➕ Новая задача...")
		updateCounter()
	})

	projectsBtn := widget.NewButtonWithIcon("Проекты", theme.FolderIcon(), func() {
		currentView = "projects"
		taskList.Refresh()
		input.Show()
		addBtn.Show()
		input.SetPlaceHolder("➕ Новый проект...")
		updateCounter()
	})

	completedBtn := widget.NewButtonWithIcon("Выполненные", theme.ConfirmIcon(), func() {
		currentView = "completed"
		taskList.Refresh()
		input.Hide()
		addBtn.Hide()
		updateCounter()
	})

	trashBtn := widget.NewButtonWithIcon("Корзина", theme.DeleteIcon(), func() {
		currentView = "trash"
		taskList.Refresh()
		input.Hide()
		addBtn.Hide()
		updateCounter()
	})

	navBar := container.NewVBox(
		inboxBtn,
		projectsBtn,
		completedBtn,
		trashBtn,
	)

	mainContent := container.NewBorder(
		container.NewVBox(
			title,
			counter,
			inputRow,
		),
		nil, nil, nil,
		container.NewScroll(taskList),
	)

	split := container.NewHSplit(
		navBar,
		mainContent,
	)
	split.Offset = 0.2

	updateCounter()
	myWindow.SetContent(split)
	myWindow.ShowAndRun()
}

// Функция сохранения данных в файл
func saveDataToFile(file string, data *Data) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}
	err = os.WriteFile(file, jsonData, 0644)
	if err != nil {
		log.Println("Error writing file:", err)
	} else {
		log.Printf("Saved: Inbox=%d, Projects=%d, Completed=%d, Trash=%d",
			len(data.Inbox), len(data.Projects), len(data.Completed), len(data.Trash))
	}
}

func showInboxActions(task *Task, data *Data, index int, window fyne.Window, onRefresh func()) {
	actions := container.NewVBox(
		widget.NewButtonWithIcon("📁 В проект", theme.FolderIcon(), func() {
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
			showProjectSelect(task, data, index, window, onRefresh)
		}),
		widget.NewButtonWithIcon("✅ Выполнено", theme.ConfirmIcon(), func() {
			task.Done = true
			data.Completed = append(data.Completed, *task)
			data.Inbox = append(data.Inbox[:index], data.Inbox[index+1:]...)
			onRefresh()
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
		}),
		widget.NewButtonWithIcon("🗑 Удалить", theme.DeleteIcon(), func() {
			data.Trash = append(data.Trash, *task)
			data.Inbox = append(data.Inbox[:index], data.Inbox[index+1:]...)
			onRefresh()
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
		}),
	)

	dialog.ShowCustom("Действия", "Закрыть", actions, window)
}

func showProjectActions(project *Task, data *Data, index int, window fyne.Window, onRefresh func()) {
	actions := container.NewVBox(
		widget.NewButtonWithIcon("📋 Подзадачи", theme.ListIcon(), func() {
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
			showSubtasks(project, window, onRefresh)
		}),
		widget.NewButtonWithIcon("✅ Выполнено", theme.ConfirmIcon(), func() {
			project.Done = true
			data.Completed = append(data.Completed, *project)
			data.Projects = append(data.Projects[:index], data.Projects[index+1:]...)
			onRefresh()
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
		}),
		widget.NewButtonWithIcon("🗑 Удалить", theme.DeleteIcon(), func() {
			data.Trash = append(data.Trash, *project)
			data.Projects = append(data.Projects[:index], data.Projects[index+1:]...)
			onRefresh()
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
		}),
	)

	dialog.ShowCustom("Действия", "Закрыть", actions, window)
}

func showCompletedActions(task *Task, data *Data, index int, window fyne.Window, onRefresh func()) {
	actions := container.NewVBox(
		widget.NewButtonWithIcon("↩ Вернуть", theme.ViewRestoreIcon(), func() {
			task.Done = false
			data.Inbox = append(data.Inbox, *task)
			data.Completed = append(data.Completed[:index], data.Completed[index+1:]...)
			onRefresh()
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
		}),
		widget.NewButtonWithIcon("🗑 Удалить", theme.DeleteIcon(), func() {
			data.Trash = append(data.Trash, *task)
			data.Completed = append(data.Completed[:index], data.Completed[index+1:]...)
			onRefresh()
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
		}),
	)

	dialog.ShowCustom("Действия", "Закрыть", actions, window)
}

func showTrashActions(task *Task, data *Data, index int, window fyne.Window, onRefresh func()) {
	actions := container.NewVBox(
		widget.NewButtonWithIcon("↩ Восстановить", theme.ViewRestoreIcon(), func() {
			if task.Done {
				data.Completed = append(data.Completed, *task)
			} else if len(task.Subtasks) > 0 {
				data.Projects = append(data.Projects, *task)
			} else {
				data.Inbox = append(data.Inbox, *task)
			}
			data.Trash = append(data.Trash[:index], data.Trash[index+1:]...)
			onRefresh()
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
		}),
		widget.NewButtonWithIcon("❌ Удалить навсегда", theme.DeleteIcon(), func() {
			data.Trash = append(data.Trash[:index], data.Trash[index+1:]...)
			onRefresh()
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
		}),
	)

	dialog.ShowCustom("Действия", "Закрыть", actions, window)
}

func showProjectSelect(task *Task, data *Data, index int, window fyne.Window, onRefresh func()) {
	var projectNames []string
	for _, p := range data.Projects {
		projectNames = append(projectNames, p.Title)
	}

	projectNames = append([]string{"✨ Создать новый проект"}, projectNames...)

	selectWidget := widget.NewSelect(projectNames, nil)

	content := container.NewVBox(
		widget.NewLabel("Выберите проект:"),
		selectWidget,
		widget.NewButtonWithIcon("Переместить", theme.MoveDownIcon(), func() {
			selected := selectWidget.Selected

			if selected == "✨ Создать новый проект" {
				newProject := Task{
					Title:    task.Title,
					Done:     false,
					Subtasks: []Task{*task},
				}
				data.Projects = append(data.Projects, newProject)
				data.Inbox = append(data.Inbox[:index], data.Inbox[index+1:]...)
				onRefresh()
				window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
			} else {
				for i, p := range data.Projects {
					if p.Title == selected {
						data.Projects[i].Subtasks = append(data.Projects[i].Subtasks, *task)
						data.Inbox = append(data.Inbox[:index], data.Inbox[index+1:]...)
						onRefresh()
						window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
						break
					}
				}
			}
		}),
	)

	dialog.ShowCustom("Переместить в проект", "Отмена", content, window)
}

func showSubtasks(project *Task, window fyne.Window, onRefresh func()) {
	content := container.NewVBox(
		widget.NewLabelWithStyle(project.Title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	if len(project.Subtasks) == 0 {
		content.Add(widget.NewLabel("Нет подзадач"))
	} else {
		for i, subtask := range project.Subtasks {
			index := i
			row := container.NewHBox(
				widget.NewButtonWithIcon("⬜", theme.ConfirmIcon(), func() {
					project.Subtasks = append(project.Subtasks[:index], project.Subtasks[index+1:]...)
					onRefresh()
					window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
					showSubtasks(project, window, onRefresh)
				}),
				widget.NewLabel(subtask.Title),
			)
			content.Add(row)
		}
	}

	entry := widget.NewEntry()
	entry.SetPlaceHolder("➕ Новая подзадача")

	addBtn := widget.NewButtonWithIcon("Добавить", theme.ContentAddIcon(), func() {
		if entry.Text != "" {
			project.Subtasks = append(project.Subtasks, Task{Title: entry.Text, Done: false})
			onRefresh()
			window.Canvas().Overlays().Remove(window.Canvas().Overlays().Top())
			showSubtasks(project, window, onRefresh)
		}
	})

	content.Add(container.NewBorder(nil, nil, nil, addBtn, entry))

	dialog.ShowCustom("Подзадачи", "Закрыть", content, window)
}
