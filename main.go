package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Task struct {
	Title string
	Done  bool
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

	// Загружаем данные
	dataDir, _ := os.UserCacheDir()
	dataDir = filepath.Join(dataDir, "gtd-data")
	os.MkdirAll(dataDir, 0755)
	dataFile := filepath.Join(dataDir, "tasks.json")

	data := &Data{}

	// Читаем существующие данные
	if file, err := os.ReadFile(dataFile); err == nil {
		json.Unmarshal(file, data)
	} else {
		// Тестовые данные
		data.Inbox = []Task{
			{Title: "Купить продукты", Done: false},
			{Title: "Сделать зарядку", Done: false},
		}
	}

	// Функция сохранения
	saveData := func() {
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		os.WriteFile(dataFile, jsonData, 0644)
	}

	currentView := "inbox"

	// Элементы интерфейса
	title := widget.NewLabelWithStyle("GTD Organizer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

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
			return widget.NewLabel("")
		},
		func(i int, o fyne.CanvasObject) {
			switch currentView {
			case "inbox":
				o.(*widget.Label).SetText(data.Inbox[i].Title)
			case "projects":
				o.(*widget.Label).SetText(data.Projects[i].Title)
			case "completed":
				o.(*widget.Label).SetText(data.Completed[i].Title)
			case "trash":
				o.(*widget.Label).SetText(data.Trash[i].Title)
			}
		},
	)

	// Обработка выбора задачи
	taskList.OnSelected = func(id widget.ListItemID) {
		switch currentView {
		case "inbox":
			showTaskActions(&data.Inbox[id], data, "inbox", id, myWindow, func() {
				saveData()
				taskList.Refresh()
			})
		}
		taskList.UnselectAll()
	}

	// Поле ввода
	input := widget.NewEntry()
	input.SetPlaceHolder("Новая задача...")
	input.OnSubmitted = func(text string) {
		if text != "" && currentView == "inbox" {
			data.Inbox = append(data.Inbox, Task{Title: text, Done: false})
			saveData()
			input.SetText("")
			taskList.Refresh()
		}
	}

	addBtn := widget.NewButton("Добавить", func() {
		if input.Text != "" && currentView == "inbox" {
			data.Inbox = append(data.Inbox, Task{Title: input.Text, Done: false})
			saveData()
			input.SetText("")
			taskList.Refresh()
		}
	})

	inputRow := container.NewBorder(nil, nil, nil, addBtn, input)

	// Кнопки навигации
	inboxBtn := widget.NewButton("Входящие", func() {
		currentView = "inbox"
		taskList.Refresh()
		input.Show()
		addBtn.Show()
	})

	projectsBtn := widget.NewButton("Проекты", func() {
		currentView = "projects"
		taskList.Refresh()
		input.Hide()
		addBtn.Hide()
	})

	completedBtn := widget.NewButton("Выполненные", func() {
		currentView = "completed"
		taskList.Refresh()
		input.Hide()
		addBtn.Hide()
	})

	trashBtn := widget.NewButton("Корзина", func() {
		currentView = "trash"
		taskList.Refresh()
		input.Hide()
		addBtn.Hide()
	})

	navBar := container.NewGridWithColumns(4, inboxBtn, projectsBtn, completedBtn, trashBtn)

	// Собираем всё вместе
	content := container.NewBorder(
		container.NewVBox(title, navBar, inputRow),
		nil, nil, nil,
		container.NewScroll(taskList),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}

func showTaskActions(task *Task, data *Data, view string, index int, window fyne.Window, onRefresh func()) {
	actions := container.NewVBox(
		widget.NewButton("✅ Выполнено", func() {
			task.Done = true
			data.Completed = append(data.Completed, *task)
			data.Inbox = append(data.Inbox[:index], data.Inbox[index+1:]...)
			onRefresh()
		}),
		widget.NewButton("🗑 Удалить", func() {
			data.Trash = append(data.Trash, *task)
			data.Inbox = append(data.Inbox[:index], data.Inbox[index+1:]...)
			onRefresh()
		}),
	)

	dialog.ShowCustom("Действия", "Закрыть", actions, window)
}
