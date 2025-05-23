package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
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

func loadTodos(db *sql.DB) ([]Todo, error) {
	var todos []Todo
	rows, err := db.Query("SELECT id, title, description, completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func getNextId(db *sql.DB) int {
	var id int
	err := db.QueryRow("SELECT MAX(id) FROM todos").Scan(&id)
	if err != nil {
		return 0
	}
	return id + 1
}

func addTodo(db *sql.DB, title, description string) error {
	todo := Todo{
		Id:          getNextId(db),
		Title:       title,
		Description: description,
		Completed:   false,
	}
	_, err := db.Exec("INSERT INTO todos (id, title, description, completed) VALUES (?, ?, ?, ?)", todo.Id, todo.Title, todo.Description, todo.Completed)
	if err != nil {
		return err
	}
	return nil
}

func completeTodo(db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE todos SET completed=true WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func removeTodo(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM todos WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

var (
	action      = flag.String("action", "add", "action (add, remove, list, complete)")
	title       = flag.String("title", "", "title of the todo")
	description = flag.String("description", "", "description of the todo")
	id          = flag.Int("id", 0, "ID of the todo to be actioned upon. Used for remove and complete")
)

func main() {
	db, err := sql.Open("sqlite3", "./todos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `CREATE TABLE if not exists todos ( id INT not null primary key, title varchar not null, description text, completed boolean);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	switch *action {
	case "add":
		handleAdd(db)
	case "remove":
		handleRemove(db)
	case "complete":
		handleComplete(db)
	case "list":
		handleList(db)
	default:
		log.Fatal("Please input a valid action.")
	}
}

func handleAdd(db *sql.DB) {
	if *title == "" || *description == "" {
		log.Fatal("Title and Description must be provided for adding a new TODO")
	}
	if err := addTodo(db, *title, *description); err != nil {
		log.Fatal(err)
	}
}

func handleRemove(db *sql.DB) {
	if *id == 0 {
		log.Fatal("ID must be provided for removing a TODO")
	}
	if err := removeTodo(db, *id); err != nil {
		log.Fatal(err)
	}
}

func handleComplete(db *sql.DB) {
	if *id == 0 {
		log.Fatal("ID must be provided for marking a Todo as complete")
	}
	if err := completeTodo(db, *id); err != nil {
		log.Fatal(err)
	}
}

func handleList(db *sql.DB) {
	todos, err := loadTodos(db)
	if err != nil {
		log.Fatal(err)
	}
	listTodos(todos)
}
