# Go Todo CLI App

A simple command-line Todo application written in Go. Todos are stored in a local JSON file.

## Features

- Add new todos with a title and description
- List all todos, grouped by incomplete and completed
- Mark todos as completed
- Remove todos by ID
- Data is persisted in `todos.json`

## Usage

Build the app:

```sh
go build -o todo
```

### Add a Todo

```sh
./todo -action=add -title="Buy milk" -description="Get 2 liters of milk"
```

### List Todos

```sh
./todo -action=list
```

### Complete a Todo

```sh
./todo -action=complete -id=1
```

### Remove a Todo

```sh
./todo -action=remove -id=1
```

## Flags

- `-action` (add, remove, list, complete) — What you want to do
- `-title` — Title for the todo (required for add)
- `-description` — Description for the todo (required for add)
- `-id` — ID of the todo (required for remove and complete)

## Example

```sh
./todo -action=add -title="Read Go book" -description="Finish chapter 3"
./todo -action=list
./todo -action=complete -id=1
./todo -action=remove -id=1
```

## Next Improvements

- Use a CLI library like [cobra](https://github.com/spf13/cobra) for better UX
- Add edit and search features
- Add unit tests
- Improve error handling
- Use a database instead of a JSON file for storage
- Add binary release for different platforms and make it available through Homebrew, Winget and Apt.
