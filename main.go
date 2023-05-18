package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
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

// soft delete
// func (t *Tasks) Remove(id uint16) {
// 	for i, task := range *t {
// 		if task.ID == id {
// 			(*t)[i].DeletedAt = time.Now().Unix()
// 		}
// 	}
// }

func (t *Tasks) list() {
	if (len(*t)) == 0 {
		fmt.Println("No tasks")
		return
	}

	for i, task := range *t {
		if task.DeletedAt != 0 {
			continue
		}
		fmt.Printf("%d - %s\n", i, task.Name)
	}
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
	case "-h":
		printHelp()
	default:
		printHelp()
	}

}
