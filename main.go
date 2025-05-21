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
	fmt.Println("- ID: ", todo.Id)
	fmt.Println("- Title: ", todo.Title)
	fmt.Println("- Description: ", todo.Description)
	fmt.Println("- Completed: ", todo.Completed)
	fmt.Println("--------")
}

func listTodos() {
	todos := getAllTodos()
	fmt.Println("TODO: \n------")
	for _, todo := range todos {
		if todo.Completed {
			continue
		}
		printTodo(todo)
	}
	fmt.Println("\nCompleted: \n------")
	for _, todo := range todos {
		if !todo.Completed {
			continue
		}
		printTodo(todo)
	}

}

func checkFile() {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		_, err := os.Create(fileName)
		if err != nil {
			log.Fatal("\n Error creating file: ", err)
		}
	}
}

func saveTodoList(todos []Todo) error {
	writeBytes, err := json.Marshal(todos)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, writeBytes, 0666)
}

func getAllTodos() []Todo {
	checkFile()
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		saveTodoList([]Todo{})
	}
	todos := []Todo{}
	if err := json.Unmarshal(fileData, &todos); err != nil {
		saveTodoList([]Todo{})
	}
	return todos
}

func completeTodo(id int) {
	todos := getAllTodos()
	for i, todo := range todos {
		if todo.Id == id {
			todo.Completed = true
			todos[i] = todo
			break
		}
	}
	saveTodoList(todos)
}

func removeTodo(id int) {
	todos := getAllTodos()
	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}
	saveTodoList(todos)
}

func getNextId() int {
	todos := getAllTodos()
	maxId := 0
	for _, todo := range todos {
		if todo.Id > maxId {
			maxId = todo.Id
		}
	}
	return maxId + 1
}

func addTodo(todo Todo) {
	todos := getAllTodos()
	todo.Id = getNextId()
	todo.Completed = false
	todos = append(todos, todo)
	writeBytes, err := json.Marshal(todos)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(fileName, writeBytes, 0666); err != nil {
		log.Fatal("There was an error writing the todo to the file", err)
	}
}

var (
	action      *string
	title       *string
	description *string
	id          *int
)

func init() {
	action = flag.String("action", "add", "action (add, remove, list, complete)")
	title = flag.String("title", "", "title of the todo")
	description = flag.String("description", "", "description of the todo")
	id = flag.Int("id", 0, "ID of the todo to be actioned upon. Used for remove and complete")

}
func main() {
	flag.Parse()

	if *action == "add" {
		if *title == "" || *description == "" {
			log.Fatal("Title and Description must be provided for adding a new TODO")
		}
		addTodo(Todo{Title: *title, Description: *description})
	} else if *action == "remove" {
		if *id == 0 {
			log.Fatal("ID must be provided for removing a TODO")
		}
		removeTodo(*id)
	} else if *action == "complete" {
		if *id == 0 {
			log.Fatal("ID must be provided for marking a Todo as complete")
		}
		completeTodo(*id)
	} else if *action == "list" {
		listTodos()
	} else {
		log.Fatal("Please input a valid action.")
	}

}
