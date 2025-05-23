# Todo CLI

A powerful command-line todo application written in Go. Manage your tasks efficiently with a simple, intuitive interface. All data is automatically stored in a SQLite database in your home directory.

## Features

- Add new todos with a title and description
- List all todos, grouped by incomplete and completed
- Mark todos as completed
- Remove todos by ID
- Automatic SQLite database management
- User-friendly command-line interface with intuitive commands
- Data stored securely in user's home directory

## Installation

### Using Go

```bash
go install github.com/AmanTahiliani/todo/cmd/todo@latest
```

### Using Make

```bash
git clone https://github.com/AmanTahiliani/todo.git
cd todo
make install
```

### Using Homebrew (macOS)

```bash
brew tap AmanTahiliani/todo
brew install todo
```

## Usage

### Add a new todo
```bash
todo add -t "Meeting" -d "Team standup at 10AM"
```

### List all todos
```bash
todo list
```

### Mark a todo as completed
```bash
todo complete -i 1
```

### Remove a todo
```bash
todo remove -i 1
```

### Get help
```bash
todo --help
```

## Data Storage

Todos are automatically stored in a SQLite database located at `~/.todo/todos.db`. The application handles all database operations transparently, so you don't need to worry about database management.

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for easier building)

### Building from source

```bash
make build
```

### Running tests

```bash
make test
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
