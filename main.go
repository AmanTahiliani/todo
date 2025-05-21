package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var fileName = "todos.json"

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func printTodo(todo Todo) {
	fmt.Printf("- ID: %d\n- Title: %s\n- Description: %s\n- Completed: %v\n--------\n",
		todo.Id, todo.Title, todo.Description, todo.Completed)
}

func listTodos(todos []Todo) {
	fmt.Println("TODO: \n------")
	for _, todo := range todos {
		if !todo.Completed {
			printTodo(todo)
		}
	}
	fmt.Println("\nCompleted: \n------")
	for _, todo := range todos {
		if todo.Completed {
			printTodo(todo)
		}
	}
}

func ensureFile() error {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	return f.Close()
}

func loadTodos() ([]Todo, error) {
	if err := ensureFile(); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return []Todo{}, nil
	}
	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return []Todo{}, nil // Reset on corrupt file
	}
	return todos, nil
}

func saveTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0666)
}

func getNextId(todos []Todo) int {
	maxId := 0
	for _, todo := range todos {
		if todo.Id > maxId {
			maxId = todo.Id
		}
	}
	return maxId + 1
}

func addTodo(title, description string) error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}
	todo := Todo{
		Id:          getNextId(todos),
		Title:       title,
		Description: description,
		Completed:   false,
	}
	todos = append(todos, todo)
	return saveTodos(todos)
}

func completeTodo(id int) error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}
	for i := range todos {
		if todos[i].Id == id {
			todos[i].Completed = true
			break
		}
	}
	return saveTodos(todos)
}

func removeTodo(id int) error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}
	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}
	return saveTodos(todos)
}

var (
	action      = flag.String("action", "add", "action (add, remove, list, complete)")
	title       = flag.String("title", "", "title of the todo")
	description = flag.String("description", "", "description of the todo")
	id          = flag.Int("id", 0, "ID of the todo to be actioned upon. Used for remove and complete")
)

func main() {
	flag.Parse()
	switch *action {
	case "add":
		handleAdd()
	case "remove":
		handleRemove()
	case "complete":
		handleComplete()
	case "list":
		handleList()
	default:
		log.Fatal("Please input a valid action.")
	}
}

func handleAdd() {
	if *title == "" || *description == "" {
		log.Fatal("Title and Description must be provided for adding a new TODO")
	}
	if err := addTodo(*title, *description); err != nil {
		log.Fatal(err)
	}
}

func handleRemove() {
	if *id == 0 {
		log.Fatal("ID must be provided for removing a TODO")
	}
	if err := removeTodo(*id); err != nil {
		log.Fatal(err)
	}
}

func handleComplete() {
	if *id == 0 {
		log.Fatal("ID must be provided for marking a Todo as complete")
	}
	if err := completeTodo(*id); err != nil {
		log.Fatal(err)
	}
}

func handleList() {
	todos, err := loadTodos()
	if err != nil {
		log.Fatal(err)
	}
	listTodos(todos)
}
