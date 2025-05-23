package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var (
	dbPath  string
	rootCmd = &cobra.Command{
		Use:   "todo",
		Short: "A simple CLI todo application",
		Long: `Todo is a CLI application for managing your tasks.
		It provides simple commands to add, list, complete, and remove todos.
		All data is stored locally in a SQLite database.`,
	}

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new todo",
		Run:   handleAdd,
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all todos",
		Run:   handleList,
	}

	completeCmd = &cobra.Command{
		Use:   "complete",
		Short: "Mark a todo as completed",
		Run:   handleComplete,
	}

	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove a todo",
		Run:   handleRemove,
	}

	title       string
	description string
	id          int
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Create .todo directory if it doesn't exist
	todoDir := filepath.Join(homeDir, ".todo")
	if err := os.MkdirAll(todoDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Set database path
	dbPath = filepath.Join(todoDir, "todos.db")

	// Add commands
	rootCmd.AddCommand(addCmd, listCmd, completeCmd, removeCmd)

	// Add flags
	addCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the todo")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the todo")
	addCmd.MarkFlagRequired("title")
	addCmd.MarkFlagRequired("description")

	completeCmd.Flags().IntVarP(&id, "id", "i", 0, "ID of the todo to complete")
	completeCmd.MarkFlagRequired("id")

	removeCmd.Flags().IntVarP(&id, "id", "i", 0, "ID of the todo to remove")
	removeCmd.MarkFlagRequired("id")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getDB() *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create table if it doesn't exist
	sqlStmt := `CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		completed BOOLEAN DEFAULT 0
	);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func handleAdd(_ *cobra.Command, _ []string) {
	db := getDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO todos (title, description) VALUES (?, ?)", title, description)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Todo added successfully!")
}

func handleList(_ *cobra.Command, _ []string) {
	db := getDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description, completed FROM todos ORDER BY completed, id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\nPending Todos:")
	fmt.Println("-------------")
	pendingFound := false

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed); err != nil {
			log.Fatal(err)
		}
		if !todo.Completed {
			pendingFound = true
			fmt.Printf("[%d] %s\n    %s\n\n", todo.Id, todo.Title, todo.Description)
		}
	}

	if !pendingFound {
		fmt.Println("No pending todos!\n")
	}

	// Reset the rows for completed todos
	rows, err = db.Query("SELECT id, title, description, completed FROM todos WHERE completed = 1 ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Completed Todos:")
	fmt.Println("---------------")
	completedFound := false

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Completed); err != nil {
			log.Fatal(err)
		}
		completedFound = true
		fmt.Printf("[%d] %s\n    %s\n\n", todo.Id, todo.Title, todo.Description)
	}

	if !completedFound {
		fmt.Println("No completed todos!")
	}
}

func handleComplete(_ *cobra.Command, _ []string) {
	db := getDB()
	defer db.Close()

	result, err := db.Exec("UPDATE todos SET completed = 1 WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected == 0 {
		fmt.Printf("No todo found with ID %d\n", id)
		return
	}

	fmt.Printf("Todo %d marked as completed!\n", id)
}

func handleRemove(_ *cobra.Command, _ []string) {
	db := getDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rowsAffected == 0 {
		fmt.Printf("No todo found with ID %d\n", id)
		return
	}

	fmt.Printf("Todo %d removed successfully!\n", id)
}
