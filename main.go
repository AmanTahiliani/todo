package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

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
	title       = flag.String("title", "", "Title of the todo item")
	description = flag.String("description", "", "Description of the todo item")
	id          = flag.Int("id", 0, "ID of the todo item to remove or mark as complete")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Todo CLI - A simple todo list manager\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [action] [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Actions:\n")
		fmt.Fprintf(os.Stderr, "  add         Add a new todo item\n")
		fmt.Fprintf(os.Stderr, "  remove      Remove a todo item\n")
		fmt.Fprintf(os.Stderr, "  complete    Mark a todo item as completed\n")
		fmt.Fprintf(os.Stderr, "  list        List all todo items\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s add -t \"Meeting\" -d \"Team standup at 10AM\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s complete -i 1\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s remove -i 2\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s list\n", os.Args[0])
	}

	flag.StringVar(title, "t", "", "shorthand for --title")
	flag.StringVar(description, "d", "", "shorthand for --description")
	flag.IntVar(id, "i", 0, "shorthand for --id")
}

func main() {
	// Check for help flag before any other operations
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		flag.Usage()
		os.Exit(0)
	}

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

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	action := os.Args[1]
	os.Args = os.Args[1:]
	flag.Parse()
	switch action {
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
