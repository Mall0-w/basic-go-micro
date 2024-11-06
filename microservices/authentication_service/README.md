TODO# Authentication Service

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Environment Setup](#environment-setup)
  - [Local Development](#local-development)
- [API Documentation](#api-documentation)
  - [Authentication Endpoints](#authentication-endpoints)
  - [Error Responses](#error-responses)

## Overview

The Authentication Service manages user authentication and session management across our microservices architecture. It handles user registration, login, password recovery, and JWT token management.

**Key Features:**
- JWT-based authentication
- Secure password hashing using bcrypt
- Email verification
- Password reset functionality
- Rate limiting
- Session management

## Getting Started

### Prerequisites
- Go 1.23 or higher
- MySQL 8.0+
- Docker & Docker Compose

### Environment Setup
1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Configure the following required environment variables:
```env
# Database Configuration
DB_HOST=your_host
DB_PORT=3306
DB_NAME=example_db
DB_USER=your_user
DB_PASSWORD=your_password

# JWT Configuration
JWT_SECRET=your_jwt_secret
PRODUCTION=false
```

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Run the service:
```bash
# Using Docker
docker-compose up -d

# Without Docker
go run main.go
```

3. Verify the service is running:
```bash
curl http://localhost:8080/auth/health
```

## API Documentation

### Authentication Endpoints

### Test Health of Service
```http
GET /auth/health
```

### Check if Authenticated
```http
GET /auth
Authorization: Bearer {auth_token} 
```

### Login
```http
POST /auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "securePassword123",
}
```

### Logout
```http
POST /auth/logout
Cookies: refresh_token: {refresh_token}
```

### Refresh Auth Token
```http
GET /auth/refresh
Cookies: refresh_token: {refresh_token}
```

<!-- For complete API documentation, see our [Swagger Documentation](http://localhost:8080/swagger/index.html) when running locally. -->

### Error Responses

| Status Code | Description |
|-------------|-------------|
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Invalid credentials |
| 403 | Forbidden - Token invalid/expired |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |