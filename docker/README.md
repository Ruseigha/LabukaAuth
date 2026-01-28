# Docker Setup for Testing

This directory contains Docker Compose configurations for running tests with containerized dependencies.

## When to Use Docker

✅ **Use Docker if:**

- You prefer isolated test environments
- You're running CI/CD pipelines
- You want to test against specific MongoDB versions
- Your laptop has no issues with Docker

❌ **Use Local MongoDB if:**

- Docker has issues on your machine
- You prefer faster test iteration
- You want to use MongoDB GUI tools easily
- You're debugging MongoDB-specific issues

## Quick Start (Docker)

```bash
# Start test MongoDB
make docker-test-up

# Run tests
make test-integration-only

# Stop MongoDB
make docker-test-down
```

## Quick Start (Local MongoDB)

See main [README.md](../README.md#local-mongodb-setup)

## Files

- `docker-compose.test.yml` - Test environment with MongoDB
