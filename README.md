# Auth Service

Production-ready authentication microservice built with Go, Clean Architecture, and MongoDB.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Setup](#setup)
  - [Option A: Local MongoDB](#option-a-local-mongodb-recommended)
  - [Option B: Docker MongoDB](#option-b-docker-mongodb)
- [Running Tests](#running-tests)
- [Development](#development)

---

## Prerequisites

- Go 1.21+

- **Choose one:**
  - **MongoDB 7.0+** (installed locally) **[Recommended for local dev]**
  - **Docker & Docker Compose** (for containerized testing)

---

## Setup

### Option A: Local MongoDB (Recommended)

**1. Install MongoDB:**

**macOS:**

```bash
brew tap mongodb/brew
brew install mongodb-community@7.0
brew services start mongodb-community@7.0
```

**Linux (Ubuntu/Debian):**

```bash
# See installation guide in docs/setup/mongodb-local.md
sudo apt-get install -y mongodb-org
sudo systemctl start mongod
```

**Windows:**

```powershell
# Download from https://www.mongodb.com/try/download/community
# Or use Chocolatey:
choco install mongodb
net start MongoDB
```

**2. Verify MongoDB:**

```bash
make mongodb-status
# Expected: ✓ MongoDB is running
```

**3. Setup project:**

```bash
# Clone repository
git clone <your-repo>
cd auth-service

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Run tests
make test-all
```

---

### Option B: Docker MongoDB

**1. Install Docker:**

- Docker Desktop: https://www.docker.com/products/docker-desktop

**2. Setup project:**

```bash
# Clone repository
git clone <your-repo>
cd auth-service

# Copy environment file
cp .env.example .env

# Install dependencies
go mod download

# Start test MongoDB
make docker-test-up

# Run tests
make test-integration

# Tests automatically start/stop Docker
```

---

## Running Tests

### Unit Tests (No Dependencies)

```bash
# Fast, no MongoDB needed
make test
```

### Integration Tests

**With Local MongoDB:**

```bash
# Make sure MongoDB is running
make mongodb-status

# Run integration tests
make test-integration

# Or clean database first
make test-integration-clean
```

**With Docker:**

```bash
# Automatically starts/stops MongoDB
make test-integration

# Or manually control:
make docker-test-up          # Start MongoDB
make test-integration-only   # Run tests
make docker-test-down        # Stop MongoDB
```

### All Tests

```bash
make test-all
```

---

## Development

### Project Structure

```
auth-service/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── config/            # Configuration
│   ├── domain/            # Domain layer (entities, value objects)
│   ├── infrastructure/    # Infrastructure (MongoDB, etc.)
│   └── delivery/          # HTTP/gRPC handlers
├── test/                  # Tests
│   ├── unit/             # Unit tests
│   └── integration/      # Integration tests
└── docker/               # Docker configurations
```

### Useful Commands

```bash
make help                  # Show all available commands
make fmt                   # Format code
make vet                   # Run go vet
make lint                  # Run linters
make test                  # Run unit tests
make test-integration      # Run integration tests
make test-all             # Run all tests
make mongodb-status        # Check MongoDB status
make mongodb-start         # Start local MongoDB
make mongodb-stop          # Stop local MongoDB
```

---

## Troubleshooting

### MongoDB Connection Issues

**Local MongoDB:**

```bash
# Check if running
make mongodb-status

# View logs (macOS)
tail -f /usr/local/var/log/mongodb/mongo.log

# View logs (Linux)
tail -f /var/log/mongodb/mongod.log

# Restart
make mongodb-stop
make mongodb-start
```

**Docker MongoDB:**

```bash
# Check container status
docker ps

# View logs
docker-compose -f docker/docker-compose.test.yml logs

# Restart
make docker-test-down
make docker-test-up
```

### Test Failures

```bash
# Clean test database
mongosh auth_service_test --eval "db.dropDatabase()"

# Or use helper
make test-integration-clean
```

---

## CI/CD

The project uses GitHub Actions with Docker for testing.

See `.github/workflows/test.yml`
