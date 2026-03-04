package main

import (
	"encoding/json"
	"fmt"
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
	myApp := app.New()
	myWindow := myApp.NewWindow("GTD Organizer")
	myWindow.Resize(fyne.NewSize(500, 600))

	// Загружаем данные
	dataDir, _ := os.UserCacheDir()
	dataDir = filepath.Join(dataDir, "gtd-data")
	os.MkdirAll(dataDir, 0755)
	dataFile := filepath.Join(dataDir, "tasks.json")

	data := &Data{}

	if file, err := os.ReadFile(dataFile); err == nil {
		json.Unmarshal(file, data)
	} else {
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
	}

	saveData := func() {
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		os.WriteFile(dataFile, jsonData, 0644)
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

	taskList.OnSelected = func(id widget.ListItemID) {
		switch currentView {
		case "inbox":
			showInboxActions(&data.Inbox[id], data, id, myWindow, func() {
				saveData()
				taskList.Refresh()
				updateCounter()
			})
		case "projects":
			showProjectActions(&data.Projects[id], data, id, myWindow, func() {
				saveData()
				taskList.Refresh()
				updateCounter()
			})
		case "completed":
			showCompletedActions(&data.Completed[id], data, id, myWindow, func() {
				saveData()
				taskList.Refresh()
				updateCounter()
			})
		case "trash":
			showTrashActions(&data.Trash[id], data, id, myWindow, func() {
				saveData()
				taskList.Refresh()
				updateCounter()
			})
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
				saveData()
				input.SetText("")
				taskList.Refresh()
				updateCounter()
			case "projects":
				data.Projects = append(data.Projects, Task{Title: text, Done: false, Subtasks: []Task{}})
				saveData()
				input.SetText("")
				taskList.Refresh()
				updateCounter()
			}
		}
	}

	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if input.Text != "" {
			switch currentView {
			case "inbox":
				data.Inbox = append(data.Inbox, Task{Title: input.Text, Done: false})
				saveData()
				input.SetText("")
				taskList.Refresh()
				updateCounter()
			case "projects":
				data.Projects = append(data.Projects, Task{Title: input.Text, Done: false, Subtasks: []Task{}})
				saveData()
				input.SetText("")
				taskList.Refresh()
				updateCounter()
			}
		}
	})

	inputRow := container.NewBorder(nil, nil, nil, addBtn, input)

	// Кнопки навигации вертикальным столбиком
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

	// Вертикальная панель навигации слева
	navBar := container.NewVBox(
		inboxBtn,
		projectsBtn,
		completedBtn,
		trashBtn,
	)

	// Основной контент справа
	mainContent := container.NewBorder(
		container.NewVBox(
			title,
			counter,
			inputRow,
		),
		nil, nil, nil,
		container.NewScroll(taskList),
	)

	// Разделитель
	split := container.NewHSplit(
		navBar,
		mainContent,
	)
	split.Offset = 0.2 // 20% на навигацию, 80% на контент

	updateCounter()
	myWindow.SetContent(split)
	myWindow.ShowAndRun()
}

// Все функции с модальными окнами, которые автоматически закрываются после действия
func showInboxActions(task *Task, data *Data, index int, window fyne.Window, onRefresh func()) {
	actions := container.NewVBox(
		widget.NewButtonWithIcon("📁 В проект", theme.FolderIcon(), func() {
			showProjectSelect(task, data, index, window, onRefresh)
		}),
		widget.NewButtonWithIcon("✅ Выполнено", theme.ConfirmIcon(), func() {
			task.Done = true
			data.Completed = append(data.Completed, *task)
			data.Inbox = append(data.Inbox[:index], data.Inbox[index+1:]...)
			onRefresh()
		}),
		widget.NewButtonWithIcon("🗑 Удалить", theme.DeleteIcon(), func() {
			data.Trash = append(data.Trash, *task)
			data.Inbox = append(data.Inbox[:index], data.Inbox[index+1:]...)
			onRefresh()
		}),
	)

	dialog.ShowCustom("Действия", "Закрыть", actions, window)
}

func showProjectActions(project *Task, data *Data, index int, window fyne.Window, onRefresh func()) {
	actions := container.NewVBox(
		widget.NewButtonWithIcon("📋 Подзадачи", theme.ListIcon(), func() {
			showSubtasks(project, window, onRefresh)
		}),
		widget.NewButtonWithIcon("✅ Выполнено", theme.ConfirmIcon(), func() {
			project.Done = true
			data.Completed = append(data.Completed, *project)
			data.Projects = append(data.Projects[:index], data.Projects[index+1:]...)
			onRefresh()
		}),
		widget.NewButtonWithIcon("🗑 Удалить", theme.DeleteIcon(), func() {
			data.Trash = append(data.Trash, *project)
			data.Projects = append(data.Projects[:index], data.Projects[index+1:]...)
			onRefresh()
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
		}),
		widget.NewButtonWithIcon("🗑 Удалить", theme.DeleteIcon(), func() {
			data.Trash = append(data.Trash, *task)
			data.Completed = append(data.Completed[:index], data.Completed[index+1:]...)
			onRefresh()
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
		}),
		widget.NewButtonWithIcon("❌ Удалить навсегда", theme.DeleteIcon(), func() {
			data.Trash = append(data.Trash[:index], data.Trash[index+1:]...)
			onRefresh()
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
			} else {
				for i, p := range data.Projects {
					if p.Title == selected {
						data.Projects[i].Subtasks = append(data.Projects[i].Subtasks, *task)
						data.Inbox = append(data.Inbox[:index], data.Inbox[index+1:]...)
						onRefresh()
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
		}
	})

	content.Add(container.NewBorder(nil, nil, nil, addBtn, entry))

	dialog.ShowCustom("Подзадачи", "Закрыть", content, window)
}
