# GC Auth Service

A robust authentication service built with Go, Gin framework, and JWT tokens. This service provides secure user authentication, registration, and token management for modern web applications.

## Features

- 🔐 **JWT-based Authentication** - Secure token-based authentication
- 👤 **User Registration & Login** - Complete user management
- 🔄 **Token Refresh** - Automatic token renewal
- 🛡️ **Password Hashing** - Secure password storage using bcrypt
- 📝 **Request Validation** - Input validation using go-playground/validator
- 🌐 **CORS Support** - Cross-origin resource sharing configuration
- 📊 **Structured Logging** - Comprehensive logging with logrus
- 🏥 **Health Check** - Service health monitoring endpoint
- 🔧 **Environment Configuration** - Flexible configuration management

## API Endpoints

### Public Endpoints
- `GET /health` - Health check endpoint
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

### Protected Endpoints (Require Authentication)
- `GET /api/v1/profile` - Get user profile
- `POST /api/v1/logout` - User logout

## Prerequisites

- Go 1.25.0 or higher
- Git

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/goldcast/gc_auth_service.git
cd gc_auth_service
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Configuration

Copy the example environment file and configure your settings:

```bash
cp env.example .env
```

Edit the `.env` file with your configuration:

```env
# Environment Configuration
ENVIRONMENT=development
PORT=8080
LOG_LEVEL=info

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY_HOURS=24

# Database Configuration (for future use)
DATABASE_URL=postgres://username:password@localhost:5432/gc_auth_db

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

### 4. Build and Run

#### Development Mode
```bash
go run main.go
```

#### Production Build
```bash
go build -o gc_auth_service main.go
./gc_auth_service
```

The service will start on port 8080 by default (or the port specified in your `.env` file).

## Usage Examples

### Health Check
```bash
curl http://localhost:8080/health
```

### User Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### User Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

### Get Profile (Protected Route)
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Refresh Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## Project Structure

```
gc_auth_service/
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Custom middleware (auth, CORS, logging, recovery)
│   ├── models/          # Data models and DTOs
│   └── routes/          # Route definitions
├── pkg/
│   ├── jwt/            # JWT token management
│   ├── logger/         # Logging utilities
│   └── password/       # Password hashing utilities
├── main.go             # Application entry point
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
├── env.example         # Environment configuration template
└── README.md           # This file
```

## Configuration

The service uses environment variables for configuration. Key configuration options:

- `ENVIRONMENT`: Application environment (development/production)
- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (debug/info/warn/error)
- `JWT_SECRET`: Secret key for JWT token signing
- `JWT_EXPIRY_HOURS`: JWT token expiration time in hours
- `DATABASE_URL`: Database connection string (for future database integration)
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed CORS origins

## Security Features

- **Password Hashing**: Uses bcrypt for secure password storage
- **JWT Tokens**: Secure token-based authentication
- **CORS Protection**: Configurable cross-origin resource sharing
- **Input Validation**: Request payload validation
- **Structured Logging**: Comprehensive audit trail

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Linting
```bash
golangci-lint run
```

## Future Enhancements

- [ ] Database integration (PostgreSQL)
- [ ] User role management
- [ ] Email verification
- [ ] Password reset functionality
- [ ] Rate limiting
- [ ] API documentation with Swagger
- [ ] Docker containerization
- [ ] Kubernetes deployment manifests

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions, please open an issue in the GitHub repository or contact the development team.
