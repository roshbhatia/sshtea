# sshtea

Manage your ssh hosts interactively with a BubbleTea-based CLI utility.

## Usage

### Dependencies

- Go `v1.16` or later
- Bubble Tea `v0.15.0` or later

### Installation

From within the repo:

1. Run `go build -o ./sshtea cmd/main/main.go`.
2. Run `./sshtea` to start the application.

### Commands

- Use arrow keys to navigate the list of hosts
- Press `a` to add a new host
- Press `e` to edit the selected host
- Press `d` to delete the selected host
- Press `h` for help
- Press `q` to quit the application

## Configuration

The application reads and writes to your SSH config file located at `~/.ssh/config`. Make sure you have the necessary permissions to read and write to this file.
