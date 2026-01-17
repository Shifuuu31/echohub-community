# EchoHub Community Forum

A modern community forum built with Go, featuring user authentication, posts, comments, likes/dislikes, and category-based organization.

## Features

- 🔐 User authentication (register, login, logout)
- 📝 Create, update, and delete posts
- 💬 Comment system with real-time updates
- 👍👎 Like/dislike functionality for posts and comments
- 🏷️ Category-based post organization
- 👤 User profiles with DiceBear avatars
- 🎨 Modern, responsive UI

## Tech Stack

- **Backend**: Go (Golang)
- **Database**: SQLite
- **Frontend**: HTML, CSS, JavaScript
- **Avatar Service**: [DiceBear API](https://www.dicebear.com/)

## Prerequisites

- Go 1.23.4 or higher
- Docker (optional, for containerized deployment)
- SQLite (included with Go)

## Quick Start

### Using Makefile (Recommended)

```bash
# Install dependencies
make install

# Run locally
make run

# Or use Docker
make start-docker      # Start Docker daemon (if using rootless Docker)
make build-docker      # Build Docker image
make run-docker-bg      # Run in background
```

### Manual Setup

```bash
# Install Go dependencies
go mod download

# Build the application
go build -o bin/echohubApp ./cmd/api/main.go

# Run the application
./bin/echohubApp
```

### Docker Setup

```bash
# Build Docker image
docker build -t echohub-community-app .

# Run container
docker run -p 8080:8080 echohub-community-app
```

## Available Make Commands

Run `make help` to see all available commands:

- `make install` - Install Go dependencies
- `make build` - Build Go application
- `make run` - Run application locally
- `make build-docker` - Build Docker image
- `make run-docker-bg` - Run Docker container in background
- `make start-docker` - Start Docker daemon (rootless)
- `make stop-docker-container` - Stop Docker container
- `make free-port` - Free port 8080 if in use
- `make clean` - Clean build artifacts

## Project Structure

```
echohub-community/
├── cmd/
│   ├── api/           # Main application entry point
│   └── web/
│       ├── handlers/  # HTTP handlers
│       └── templates/ # HTML templates
├── internal/
│   ├── database/      # Database files and migrations
│   └── models/       # Data models
├── assets/           # Static assets (CSS, JS, images)
├── Makefile          # Build and run commands
└── Dockerfile        # Docker configuration
```

## Configuration

The application runs on port `8080` by default. You can change this by setting the `PORT` environment variable:

```bash
export PORT=3000
make run
```

## Database

The application uses SQLite database located at `internal/database/echohub-community.db`. The database is automatically created on first run.

## Avatar Generation

User avatars are generated using [DiceBear API](https://www.dicebear.com/) with the "adventurer" style. Avatars are deterministic based on username and gender, ensuring consistent avatars for each user.

## Development

```bash
# Run in development mode (with auto-reload if air/realize is installed)
make dev

# Run tests
make test
```

## License

[Add your license here]

## Links

- [UI Design on Figma](https://www.figma.com/design/QCQgn3VWW7m4NX0lEDCfNC/Figma-Design-Forum-Concept?node-id=0-1&t=UmWF9mEliRamw4b9-1)
- [GitHub Repository](https://github.com/Shifuuu31/forum)
