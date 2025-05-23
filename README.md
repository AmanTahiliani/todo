# Go Todo CLI App

A simple command-line Todo application written in Go. Todos are stored in a SQLite database.

## Features

- Add new todos with a title and description
- List all todos, grouped by incomplete and completed
- Mark todos as completed
- Remove todos by ID
- Data is persisted in SQLite database
- User-friendly command-line interface with action commands

## Dependencies

- Go 1.x
- github.com/mattn/go-sqlite3

## Usage

Build the app:

```sh
go build -o todo
```

### Command Structure

The app uses a command-based interface where the action comes first, followed by flags:

```sh
./todo [action] [flags]
```

Available actions:
- `add` - Add a new todo item
- `remove` - Remove a todo item
- `complete` - Mark a todo item as completed
- `list` - List all todo items

### Flags

- `-t, --title` - Title for the todo (required for add)
- `-d, --description` - Description for the todo (required for add)
- `-i, --id` - ID of the todo (required for remove and complete)
- `-h, --help` - Display help information

### Examples

Add a new todo:
```sh
./todo add -t "Buy milk" -d "Get 2 liters of milk"
```

List all todos:
```sh
./todo list
```

Complete a todo:
```sh
./todo complete -i 1
```

Remove a todo:
```sh
./todo remove -i 2
```

Display help:
```sh
./todo -h
```

## Next Improvements

- Use a CLI library like [cobra](https://github.com/spf13/cobra) for better UX
- Add edit and search features
- Add unit tests
- Improve error handling
- Add binary release for different platforms and make it available through Homebrew, Winget and Apt
- Add due dates and priority levels for todos
- Implement categories/tags for better organization
- Add data export/import functionality
