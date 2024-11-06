# Microservices Architecture
- [Getting Started](#getting-started)
  - [Testing](#testing)
  - [Running the Backend](#running-the-backend)
  - [Prerequisites](#prerequisites)
- [Services](#services)
  - [Gateway](#gateway)
  - [Authenticate Service](#authenticate-service)
  - [User Service](#user-service)
- [Microservice Structure](#microservice-structure)
  - [Core Components](#core-components)
  - [Additional Information](#additional-information)

## Getting Started

### Testing
```bash
go test ./tests
```

### Running the Backend
```bash
docker-compose up -d --build
```

### Prerequisites
Please ensure you have all the private files for each service. This will be more manageable after integrating Kubernetes.

## Services

### Gateway
This is the API gateway for the microservice architecture. Currently written in Go until we implement an automated gateway solution with Kubernetes. When adding a service, make sure to also add it to the gateway.

For detailed gateway configuration, see [gateway/README.md](./gateway/README.md)

### Authenticate Service
API service that handles all authentication of users.

For authentication flow details, see [authentication_service/README.md](./microservices/authentication_service/README.md)

### User Service
API service that handles all user profiles.

For user service documentation, see [user_service/README.md](./microservices/user_service/README.md)

## Microservice Structure
Each microservice follows the MVC (Model-View-Controller) pattern. Below are the key components:

### Core Components

| Component | Description |
|-----------|-------------|
| `main.go` | Entry point for the service, initializes Gin and injects dependencies |
| `.env` | Environment configuration file |
| `Dockerfile` | Standardized container definition |
| `config/` | Environment variables and service configurations |
| `controller/` | MVC controllers defining routes and request handling |
| `dtos/` | Data Transfer Objects for standardized JSON handling |
| `errors/` | Custom error types and response formatting |
| `models/` | Database schema definitions for GORM |
| `repository/` | Database interaction layer with dependency injection |
| `service/` | Business logic implementation |

### Additional Information
For detailed implementation guides:
- [TODO List](./docs/TODO.md)