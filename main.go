package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

type STATUS uint8

const (
	TODO STATUS = iota
	DOING
	ONTEST
	DONE
)

const TASK_FILENAME = "tasks.json"

type Task struct {
	Name      string `json:"label"`
	Status    STATUS `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

type Tasks map[uint16]Task

func saveTasks(filename string, t Tasks) {
	if len(t) == 0 || filename == "" {
		fmt.Println("Filename and tasks are required")
		return
	}
	tasksJSON, _ := json.MarshalIndent(t, "", "  ")
	os.WriteFile(filename, tasksJSON, 0644)
}

func (t *Tasks) add(name string) {
	if name == "" {
		fmt.Println("Task name is required")
		os.Exit(1)
	}

	task := Task{
		Name:      name,
		Status:    TODO,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	(*t)[uint16(len(*t)+1)] = task

	saveTasks(TASK_FILENAME, *t)

	fmt.Printf("Task \"%s\" added \n", name)
}

func (t *Tasks) remove(id uint16) {
	value, exists := (*t)[id]
	if exists {
		value.DeletedAt = time.Now().Unix()
		(*t)[id] = value
		saveTasks(TASK_FILENAME, *t)
		fmt.Printf("Task \"%s\" removed \n", value.Name)
		return
	}
	fmt.Println("Task not found")
}

func getTaskStatus(status STATUS) string {
	switch status {
	case TODO:
		return "TODO"
	case DOING:
		return "DOING"
	case ONTEST:
		return "ONTEST"
	case DONE:
		return "DONE"
	default:
		return "UNKNOWN"
	}
}

func unixToTime(unix int64) string {
	return time.Unix(unix, 0).Local().Format("2006-01-02 15:04:05")
}

func (t *Tasks) list() {
	if (len(*t)) == 0 {
		fmt.Println("No tasks")
		return
	}

	taskRows := make([]table.Row, 100)

	for i, task := range *t {
		if task.DeletedAt != 0 {
			continue
		}
		taskRows = append(taskRows, table.Row{i, task.Name, getTaskStatus(task.Status), unixToTime(task.CreatedAt), unixToTime(task.UpdatedAt)})
	}

	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.AppendHeader(table.Row{"ID", "Task", "Status", "Created At", "Updated At"})
	tab.AppendRows(taskRows)
	tab.SortBy([]table.SortBy{
		{Name: "ID", Mode: table.Asc},
	})
	tab.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Task", WidthMax: 40},
	})
	tab.Render()
}

func printHelp() {
	fmt.Println("Usage: ")
	fmt.Println("  -a <task> \t - add new task")
	fmt.Println("  -d <id> \t - delete a task")
	fmt.Println("  -l \t\t - show the tasks")
	fmt.Println("  -h \t\t - help")
	os.Exit(0)
}

func main() {
	// Create tasks
	tasks := make(Tasks)

	// Read tasks from file
	tasksJSON, err := os.ReadFile(TASK_FILENAME)
	if err != nil {
		fmt.Println("No tasks")
	} else {
		json.Unmarshal(tasksJSON, &tasks)
	}

	if (len(os.Args)) < 2 {
		printHelp()
	}

	switch os.Args[1] {
	case "-a":
		tasks.add(strings.Join(os.Args[2:], " "))
	case "-l":
		tasks.list()
	case "-d":
		tasks.remove(uint16(os.Args[2][0] - '0'))
	case "-h":
		printHelp()
	default:
		printHelp()
	}

}
