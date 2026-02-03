# Auth Service

Production-ready authentication microservice built with Go, Clean Architecture, and MongoDB.

[![CI/CD](https://github.com/your-org/auth-service/workflows/CI%2FCD%20Pipeline/badge.svg)](https://github.com/your-org/auth-service/actions)
[![codecov](https://codecov.io/gh/your-org/auth-service/branch/main/graph/badge.svg)](https://codecov.io/gh/your-org/auth-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-org/auth-service)](https://goreportcard.com/report/github.com/your-org/auth-service)

## 🚀 Features

- **Clean Architecture** - Domain-driven design with clear separation of concerns
- **Dual Transport** - HTTP REST API + gRPC for service-to-service
- **JWT Authentication** - Secure token-based authentication
- **MongoDB** - NoSQL database with proper indexing
- **Production Ready** - Docker, CI/CD, monitoring, graceful shutdown
- **Comprehensive Testing** - Unit tests, integration tests, 95%+ coverage
- **Security** - Rate limiting, password hashing (bcrypt), input validation

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Configuration](#configuration)
- [Contributing](#contributing)

## 🔧 Prerequisites

- **Go** 1.21+
- **MongoDB** 7.0+
- **Docker** & Docker Compose (optional)
- **Protocol Buffers** compiler (for gRPC)

## ⚡ Quick Start

### Local Development
```bash
# Clone repository
git clone https://github.com/your-org/auth-service.git
cd auth-service

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Start MongoDB (if not running)
make mongodb-start

# Run server
make run
```

Server will start on:
- HTTP: `http://localhost:8001`
- gRPC: `localhost:9001`

### Docker
```bash
# Start with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

## 🏛️ Architecture
```
auth-service/
├── cmd/server/              # Application entry point
├── internal/
│   ├── config/             # Configuration
│   ├── domain/             # Business logic (entities, value objects)
│   ├── usecase/            # Use cases (application logic)
│   ├── infrastructure/     # External dependencies (DB, security)
│   └── delivery/           # Transport layer (HTTP, gRPC)
├── test/
│   ├── unit/              # Unit tests
│   └── integration/       # Integration tests
├── proto/                 # Protocol buffer definitions
└── docker/                # Docker configs
```

### Clean Architecture Layers
```
┌─────────────────────────────────────────┐
│          Delivery (HTTP/gRPC)           │ ← External interface
├─────────────────────────────────────────┤
│          Use Cases (Business)           │ ← Application logic
├─────────────────────────────────────────┤
│     Domain (Entities, Value Objects)    │ ← Core business rules
└─────────────────────────────────────────┘
```

## 📡 API Documentation

### HTTP Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/signup` | Register new user |
| POST | `/api/v1/auth/login` | Authenticate user |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/api/v1/auth/validate` | Validate token (protected) |
| GET | `/health` | Health check |

### gRPC Services
```protobuf
service AuthService {
  rpc Signup(SignupRequest) returns (AuthResponse);
  rpc Login(LoginRequest) returns (AuthResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (AuthResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}
```

### Example: Signup

**HTTP:**
```bash
curl -X POST http://localhost:8001/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecureP@ss123"
  }'
```

**gRPC:**
```bash
grpcurl -plaintext \
  -d '{"email":"user@example.com","password":"SecureP@ss123"}' \
  localhost:9001 auth.AuthService/Signup
```

**Response:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 🛠️ Development

### Available Commands
```bash
make help                  # Show all commands
make run                   # Run server
make test                  # Run unit tests
make test-integration      # Run integration tests
make test-all             # Run all tests
make lint                  # Run linters
make proto                 # Generate protobuf code
make build                 # Build binary
make docker-build          # Build Docker image
```

### Project Structure Explained

- **cmd/server/** - Application entry point
- **internal/domain/** - Core business logic (NO external dependencies)
- **internal/usecase/** - Application business logic (orchestrates domain)
- **internal/infrastructure/** - External integrations (MongoDB, JWT)
- **internal/delivery/** - Transport layers (HTTP, gRPC)

## 🧪 Testing
```bash
# Unit tests (fast, no dependencies)
make test

# Integration tests (requires MongoDB)
make test-integration

# All tests
make test-all

# With coverage
make test-coverage

# Benchmarks
make benchmark
```

### Test Coverage

- **Domain Layer**: 95%+
- **Use Case Layer**: 95%+
- **Overall**: 90%+

## 🚢 Deployment

### Docker
```bash
# Build image
docker build -t auth-service:latest .

# Run container
docker run -p 8001:8001 -p 9001:9001 \
  -e MONGO_URI=mongodb://mongo:27017 \
  -e JWT_SECRET_KEY=your-secret \
  auth-service:latest
```

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    metadata:
      labels:
        app: auth-service
    spec:
      containers:
      - name: auth-service
        image: ghcr.io/your-org/auth-service:latest
        ports:
        - containerPort: 8001
        - containerPort: 9001
        env:
        - name: MONGO_URI
          valueFrom:
            secretKeyRef:
              name: mongodb-credentials
              key: uri
```

## ⚙️ Configuration

Configuration via environment variables:
```bash
# Application
APP_ENV=production
APP_VERSION=1.0.0

# Server
HTTP_PORT=8001
GRPC_PORT=9001

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=auth_service

# JWT
JWT_SECRET_KEY=your-secret-min-32-chars
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=168h
```

See `.env.example` for complete configuration.

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### Code Standards

- Follow Go best practices
- Write tests for new features
- Update documentation
- Run `make lint` before committing

## 📄 License

MIT License - see [LICENSE](LICENSE) file

## 👥 Authors

- Your Name - [@yourhandle](https://github.com/yourhandle)

## 🙏 Acknowledgments

- Clean Architecture by Robert C. Martin
- Domain-Driven Design by Eric Evans
- Go microservices community
