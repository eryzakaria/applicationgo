# SuiteMedia API

A production-ready REST API built with Go (Golang) and Gin framework, featuring JWT authentication, PostgreSQL database, Redis caching, and comprehensive middleware stack.

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- Redis 6+ (optional)
- Docker (optional, for running services)

### Installation

1. **Clone the repository**
```bash
cd d:\Works\SuiteMedia
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup environment variables**
```bash
cp .env.example .env
```

Edit `.env` file with your configuration:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=suitemedia
```

### Running with Docker (Recommended)

1. **Start PostgreSQL and Redis**
```bash
# PostgreSQL
docker run -d --name postgres \
  -p 5432:5432 \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=suitemedia \
  postgres:15-alpine

# Redis (optional)
docker run -d --name redis \
  -p 6379:6379 \
  redis:7-alpine
```

2. **Run the application**
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:3000`

### Running without Docker

1. **Install PostgreSQL locally**
   - Windows: Download from [postgresql.org](https://www.postgresql.org/download/windows/)
   - Linux: `sudo apt-get install postgresql`
   - macOS: `brew install postgresql`

2. **Create database**
```sql
CREATE DATABASE suitemedia;
CREATE USER myuser WITH PASSWORD 'mypassword';
GRANT ALL PRIVILEGES ON DATABASE suitemedia TO myuser;
```

3. **Run the application**
```bash
go run cmd/api/main.go
```

## 📝 API Endpoints

### Health Check
```bash
# Check if service is running
curl http://localhost:3000/health

# Check if service is ready (DB + Redis)
curl http://localhost:3000/ready
```

### Authentication

**Register a new user:**
```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

**Refresh Token:**
```bash
curl -X POST http://localhost:3000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your_refresh_token"
  }'
```

### User Management

**Get all users (requires authentication):**
```bash
curl http://localhost:3000/api/v1/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get user by ID:**
```bash
curl http://localhost:3000/api/v1/users/{id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get current user profile:**
```bash
curl http://localhost:3000/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Update profile:**
```bash
curl -X PUT http://localhost:3000/api/v1/users/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith"
  }'
```

**Create user (admin only):**
```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "password123",
    "first_name": "Alice",
    "last_name": "Johnson",
    "role": "user"
  }'
```

**Update user (admin only):**
```bash
curl -X PUT http://localhost:3000/api/v1/users/{id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Updated Name",
    "is_active": true
  }'
```

**Delete user (admin only):**
```bash
curl -X DELETE http://localhost:3000/api/v1/users/{id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Products

**List products:**
```bash
curl http://localhost:3000/api/v1/products?page=1&limit=10 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Get product by ID:**
```bash
curl http://localhost:3000/api/v1/products/{id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 🏗️ Project Structure

```
SuiteMedia/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── config/
│   └── config.go                # Configuration management
├── internal/
│   ├── database/
│   │   └── connection.go        # Database connection & migrations
│   ├── handlers/
│   │   ├── auth_handler.go      # Authentication endpoints
│   │   ├── user_handler.go      # User CRUD endpoints
│   │   ├── product_handler.go   # Product endpoints
│   │   └── health_handler.go    # Health check endpoints
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication
│   │   ├── cors.go              # CORS configuration
│   │   ├── logger.go            # Request logging
│   │   └── metrics.go           # Prometheus metrics
│   ├── models/
│   │   ├── user.go              # User models & DTOs
│   │   └── product.go           # Product models & DTOs
│   ├── repository/
│   │   ├── user_repository.go   # User data access
│   │   └── product_repository.go
│   └── service/
│       ├── auth_service.go      # Auth business logic
│       ├── user_service.go      # User business logic
│       └── product_service.go
├── pkg/
│   ├── logger/
│   │   └── logger.go            # Logging utility
│   ├── redis/
│   │   └── redis.go             # Redis client
│   └── response/
│       └── response.go          # API response helpers
├── .env                         # Environment variables
├── .env.example                 # Environment template
├── go.mod                       # Go module dependencies
└── README.md                    # This file
```

## 🔧 Configuration

All configuration is done via environment variables in `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | API server port | 8080 |
| `SERVER_ENV` | Environment (development/production) | development |
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `DB_USER` | Database user | postgres |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | suitemedia |
| `REDIS_HOST` | Redis host | localhost |
| `REDIS_PORT` | Redis port | 6379 |
| `JWT_SECRET` | JWT signing secret | - |
| `JWT_EXPIRATION_HOURS` | Access token expiration | 24 |
| `JWT_REFRESH_EXPIRATION_DAYS` | Refresh token expiration | 30 |

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication:

1. Register or login to get an access token
2. Include the token in the `Authorization` header: `Bearer YOUR_TOKEN`
3. Access tokens expire after 24 hours (configurable)
4. Use the refresh token to get a new access token

### Roles

- **user**: Regular user with read access
- **admin**: Full access to all endpoints

## 🐳 Docker

**Build Docker image:**
```bash
docker build -t suitemedia-api .
```

**Run with Docker Compose:**
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DB_HOST=postgres
      - DB_PASSWORD=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: suitemedia
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

## 📊 Monitoring

**Prometheus Metrics:**
```bash
curl http://localhost:3000/metrics
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestName ./internal/service
```

## 🚀 Deployment

### Deploy to Kubernetes (Helm)

```bash
# Install with Helm
cd helm/suitemedia

# Development
helm install suitemedia . -f values-staging.yaml

# Production
helm install suitemedia . -f values-production.yaml
```

### Environment-specific Configurations

- **Staging**: `values-staging.yaml` - 2-5 replicas, spot instances
- **Production**: `values-production.yaml` - 6-20 replicas, on-demand instances

## 📚 Documentation

Additional documentation available:

- [Infrastructure Documentation](INFRASTRUCTURE_DOCUMENTATION.md)
- [Kubernetes Guide](K8S_INFRASTRUCTURE_GUIDE.md)
- [CI/CD Documentation](CI-CD-DOCUMENTATION.md)
- [Helm Chart README](helm/suitemedia/README.md)

## 🛠️ Development

**Run with hot reload (using Air):**
```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

**Format code:**
```bash
go fmt ./...
```

**Lint code:**
```bash
golangci-lint run
```

## 🐛 Troubleshooting

**Database connection refused:**
- Ensure PostgreSQL is running: `docker ps` or check service status
- Verify credentials in `.env` file
- Check if database exists: `psql -U myuser -d suitemedia`

**Redis connection error:**
- Redis is optional - the app will still run
- Start Redis: `docker start redis` or install locally
- Disable Redis by commenting out Redis initialization in `main.go`

**Port already in use:**
- Change `SERVER_PORT` in `.env` file
- Or kill the process: `lsof -ti:3000 | xargs kill` (macOS/Linux)

## 📄 License

MIT License - see LICENSE file for details

## 👥 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature-name`
3. Commit changes: `git commit -am 'Add feature'`
4. Push to branch: `git push origin feature-name`
5. Submit a Pull Request

## 📧 Support

For issues and questions, please open an issue on GitHub.
