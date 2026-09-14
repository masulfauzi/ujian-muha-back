# Fiber Backend API

Production-ready backend API built with Go, Fiber, and PostgreSQL following clean architecture principles.

## Features

- ✅ Modern REST API with Fiber v2
- ✅ PostgreSQL database with GORM ORM
- ✅ JWT Authentication & Authorization
- ✅ Clean Architecture & Repository Pattern
- ✅ Modular Structure
- ✅ Password hashing with bcrypt
- ✅ Request validation
- ✅ Global error handling
- ✅ Logger middleware
- ✅ CORS support
- ✅ Auto database migration

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator
- **Environment**: godotenv

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── configs/
│   ├── app.go                   # App configuration
│   ├── database.go              # Database configuration
│   └── jwt.go                   # JWT configuration
├── internal/
│   ├── modules/
│   │   ├── auth/
│   │   │   ├── controller/
│   │   │   ├── service/
│   │   │   ├── repository/
│   │   │   ├── dto/
│   │   │   ├── model/
│   │   │   ├── routes/
│   │   │   └── validator/
│   │   └── user/
│   │       ├── controller/
│   │       ├── service/
│   │       ├── repository/
│   │       ├── dto/
│   │       ├── model/
│   │       ├── routes/
│   │       └── validator/
│   ├── middleware/
│   │   ├── jwt_auth.go
│   │   ├── logger.go
│   │   ├── cors.go
│   │   └── recovery.go
│   ├── database/
│   │   ├── connection.go
│   │   └── migrate.go
│   ├── helpers/
│   │   └── response.go
│   ├── utils/
│   │   └── password.go
│   └── constants/
│       └── constants.go
├── migrations/                  # Database migrations (if needed)
├── go.mod
├── go.sum
├── .env
├── .env.example
└── README.md
```

## Installation

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

### Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd backend
```

2. **Install dependencies**
```bash
go mod download
go mod tidy
```

3. **Configure environment variables**
```bash
cp .env.example .env
```

Edit `.env` file with your configuration:

```env
APP_NAME=Fiber Backend API
APP_PORT=3000
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=fiber_backend

JWT_SECRET=supersecretkey
JWT_EXPIRED=72h

LOG_LEVEL=info
```

4. **Create database**
```bash
createdb fiber_backend
```

Or using PostgreSQL client:

```sql
CREATE DATABASE fiber_backend;
```

## Running the Server

### Development Mode

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:3000`

### Production Build

```bash
go build -o backend cmd/server/main.go
./backend
```

## API Endpoints

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "Fiber Backend API"
}
```

### Authentication

#### Register
```http
POST /api/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Register successfully",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "token": "jwt_token",
    "role": "user"
  }
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successfully",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "token": "jwt_token",
    "role": "user"
  }
}
```

#### Get Current User
```http
GET /api/auth/me
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Get current user successfully",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
  }
}
```

### Users

#### Get All Users
```http
GET /api/users
```

**Response:**
```json
{
  "success": true,
  "message": "Get all users successfully",
  "data": [
    {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "status": "active",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### Get User By ID
```http
GET /api/users/:id
```

**Response:**
```json
{
  "success": true,
  "message": "Get user successfully",
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "status": "active",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### Create User (Authenticated)
```http
POST /api/users
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": "uuid",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "role": "user",
    "status": "active",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### Update User (Authenticated)
```http
PUT /api/users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Jane Doe Updated",
  "email": "jane-updated@example.com",
  "role": "admin"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": "uuid",
    "name": "Jane Doe Updated",
    "email": "jane-updated@example.com",
    "role": "admin",
    "status": "active",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### Delete User (Authenticated)
```http
DELETE /api/users/:id
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": null
}
```

## Validation Rules

### Register & Create User
- **name**: Required
- **email**: Required, must be valid email format
- **password**: Required, minimum 6 characters

### Login
- **email**: Required, must be valid email format
- **password**: Required

### Update User
- **name**: Optional
- **email**: Optional, must be valid email format if provided
- **role**: Optional

## Authentication

All protected endpoints require JWT token in Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

The JWT token contains:
- `user_id`: User UUID
- `email`: User email
- `role`: User role
- `exp`: Token expiration time (default 72 hours)

## Error Responses

### Validation Error
```json
{
  "success": false,
  "message": "Validation error",
  "errors": null
}
```

### Unauthorized
```json
{
  "success": false,
  "message": "Missing authorization header",
  "errors": null
}
```

### Not Found
```json
{
  "success": false,
  "message": "Resource not found",
  "errors": null
}
```

### Internal Server Error
```json
{
  "success": false,
  "message": "Internal server error",
  "errors": null
}
```

## Development

### Code Structure

- **Controllers**: Handle HTTP requests/responses
- **Services**: Business logic
- **Repositories**: Data access layer
- **Models**: Database models
- **DTOs**: Data transfer objects (request/response)
- **Middleware**: Request/response interceptors
- **Validators**: Input validation

### Adding New Module

1. Create folder under `internal/modules/<module_name>`
2. Create subfolders: `controller`, `service`, `repository`, `model`, `dto`, `routes`, `validator`
3. Implement interfaces following existing pattern
4. Create routes file
5. Register routes in `cmd/server/main.go`

### Best Practices

- Use dependency injection
- Keep business logic in services
- Use interfaces for flexibility
- Validate at boundaries
- Use constants for magic strings
- Log important operations
- Handle errors gracefully
- Use meaningful error messages

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | Fiber Backend API |
| `APP_PORT` | Server port | 3000 |
| `APP_ENV` | Environment (development/production) | development |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database user | postgres |
| `DB_PASSWORD` | Database password | postgres |
| `DB_NAME` | Database name | fiber_backend |
| `JWT_SECRET` | JWT secret key | supersecretkey |
| `JWT_EXPIRED` | JWT expiration duration | 72h |
| `LOG_LEVEL` | Log level | info |

## Troubleshooting

### Database Connection Error

**Error**: `failed to connect to database`

**Solution**:
1. Verify PostgreSQL is running
2. Check database credentials in `.env`
3. Ensure database exists: `createdb fiber_backend`

### JWT Token Invalid

**Error**: `Invalid token`

**Solution**:
1. Verify token format: `Bearer <token>`
2. Check JWT_SECRET matches
3. Verify token hasn't expired

### CORS Error

**Error**: `Cross-Origin Request Blocked`

**Solution**:
- CORS is enabled for all origins by default
- Modify CORS config in `internal/middleware/cors.go` if needed

## Performance Tips

- Use connection pooling for database
- Enable caching where applicable
- Monitor query performance
- Use indexes on frequently queried columns
- Implement pagination for large datasets
- Use read replicas for scaling reads

## Security

- Never commit `.env` file
- Use strong JWT_SECRET in production
- Validate all inputs
- Hash passwords with bcrypt
- Use HTTPS in production
- Implement rate limiting
- Sanitize database queries (GORM handles this)
- Regular security audits

## License

MIT License - see LICENSE file for details

## Support

For issues and questions, please create an issue in the repository.

## Contributors

- Your Name (your.email@example.com)

---

**Made with ❤️ using Golang & Fiber**
