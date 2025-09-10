# XArch - Go Project Base Code

A ready-to-use Go project template built with clean architecture principles, supporting both HTTP REST API and gRPC protocols. Perfect as a starting point for any Go project - from simple web applications to complex microservices.

## 🚀 Features

- **Ready-to-Use**: Complete project structure with all essential components
- **Dual Protocol Support**: HTTP REST API and gRPC services (use what you need)
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Multi-Database Support**: PostgreSQL, MySQL, and SQLite
- **Additional Storage**: MongoDB and Redis support
- **Database Migrations**: Automated database schema management
- **Seeding**: Database seeding capabilities
- **Internationalization**: Multi-language support (English & Indonesian)
- **Structured Logging**: Comprehensive logging with multiple output modes
- **Response Management**: Centralized response handling and error management
- **Health Monitoring**: Built-in health check endpoints
- **Configuration Management**: Environment-based configuration
- **Validation**: Input validation with custom validators
- **Security**: Password hashing with bcrypt
- **Flexible**: Easily customizable for any Go project type

## 🏗️ Architecture

```
├── cmd/                    # Application entry points
├── config/                 # Configuration management
├── constant/               # Application constants
├── database/               # Database migrations and seeders
├── domain/                 # Domain layer (models, repositories, services)
├── driver/                 # External service drivers
├── infrastructure/         # Infrastructure layer
│   ├── grpc/               # gRPC server implementation
│   └── http/               # HTTP server implementation
├── utils/                  # Utility functions
└── main.go                 # Main application entry point
```

## 🛠️ Tech Stack

- **Language**: Go 1.24.6
- **Web Framework**: Echo v4
- **gRPC**: Google gRPC
- **Database ORM**: GORM
- **Databases**: PostgreSQL, MySQL, SQLite
- **NoSQL**: MongoDB
- **Cache**: Redis
- **Logging**: Custom logger with file/terminal output
- **Validation**: go-playground/validator
- **Configuration**: Environment variables
- **Migration**: Custom migration system
- **Seeding**: Custom seeder system

## 📋 Prerequisites

- Go 1.24.6 or higher
- PostgreSQL/MySQL/SQLite (choose one)
- MongoDB (optional)
- Redis (optional)

## 🚀 Quick Start

### 1. Installation

#### Option 1: Manual Installation
```bash
git clone https://github.com/risoftinc/xarch.git
cd xarch
```

#### Option 2: Using Elsa CLI Tool (Recommended)
First, install the Elsa CLI tool:
```bash
go install github.com/risoftinc/elsa/cmd/elsa@latest
```

Then create a new project using XArch template:
```bash
elsa new xarch[@version] <project-name> --module=<your-module>
```

**Example:**
```bash
elsa new xarch@latest my-project --module=risoftinc.com/my-project
```

This will create a new project with all the base code structure ready for your Go application.

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Configuration

Create a `.env` file in the root directory:

```env
# Server Configuration
SERVER=localhost
PORT=9000
USING_SECURE=false

# gRPC Configuration
GRPC_SERVER=localhost
GRPC_PORT=9001

# Database Configuration
DB_TYPE=postgres
DB_USER=root
DB_PASS=password
DB_SERVER=localhost
DB_PORT=5432
DB_NAME=xarch
DB_TIME_ZONE=Asia/Jakarta
DB_SSL_MODE=disable

# Database Connection Pool
DB_MAX_IDLE_CON=10
DB_MAX_OPEN_CON=100
DB_MAX_LIFE_TIME=10
DB_DEBUG=false

# MongoDB Configuration (Optional)
MONGODB_URI=mongodb://localhost:27017
MONGODB_USERNAME=
MONGODB_PASSWORD=
MONGODB_DATABASE=xarch
MONGODB_MAX_POOL_SIZE=100
MONGODB_MIN_POOL_SIZE=5
MONGODB_MAX_IDLE_TIME=30m
MONGODB_TIMEOUT=10s

# Redis Configuration (Optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_USERNAME=root
REDIS_PASSWORD=
REDIS_DB=0
REDIS_MAX_RETRIES=3
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s
REDIS_IDLE_TIMEOUT=5m

# Logger Configuration
LOG_OUTPUT_MODE=both
LOG_LEVEL=debug
LOG_DIR=logger

# Response Manager Configuration
RESPONSE_MANAGER_METHOD=file
RESPONSE_MANAGER_PATH=config/config.json
RESPONSE_MANAGER_INTERVAL=5m
```

### 4. Database Setup

#### For PostgreSQL:
```bash
createdb xarch
```

#### For MySQL:
```sql
CREATE DATABASE xarch;
```

#### For SQLite:
No setup required, database file will be created automatically.

### 5. Run Database Migrations

```bash
# Run migrations
go run cmd/seeder/main.go
```

### 6. Start the Application

```bash
go run main.go
```

The application will start both HTTP and gRPC servers:
- HTTP Server: `http://localhost:9000`
- gRPC Server: `localhost:9001`

**Note:** You can disable either server by modifying the configuration or removing the respective server initialization code if you only need one protocol.

## 📡 API Endpoints

### HTTP REST API

#### Health Check
```http
GET /health/metric
```

**Response:**
```json
{
  "status": 200,
  "message": "Data retrieved successfully",
  "data": {
    "status": {
      "database": "healthy",
      "redis": "healthy",
      "mongodb": "healthy"
    },
    "database": {
      "MaxOpenConnections": 100,
      "OpenConnections": 5,
      "InUse": 2,
      "Idle": 3,
      "WaitCount": 0,
      "WaitDuration": 0,
      "MaxIdleClosed": 0,
      "MaxIdleTimeClosed": 0,
      "MaxLifetimeClosed": 0
    }
  }
}
```

### gRPC API

#### Health Service
```protobuf
service HealthService {
  rpc GetHealthMetric(HealthMetricRequest) returns (HealthMetricResponse);
}
```

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(255) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `roles` VARCHAR(255) NOT NULL,
  `salary` DOUBLE NOT NULL,
  `created_by` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` BIGINT UNSIGNED NULL,
  `updated_at` DATETIME NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 🔧 Configuration

The application supports multiple configuration methods:

### Environment Variables
All configuration is managed through environment variables with sensible defaults.

### Response Management
Centralized response management with support for:
- Multiple languages (English, Indonesian)
- Custom message templates
- HTTP and gRPC status code mapping
- Dynamic configuration reloading

### Database Support
- **PostgreSQL**: Full support with SSL configuration
- **MySQL**: Full support with charset and timezone configuration
- **SQLite**: File-based database for development

## 🧪 Testing

Run tests for specific packages:

```bash
# Test bcrypt utilities
go test ./utils/bcrypt/

# Run all tests
go test ./...
```

## 📝 Logging

The application includes comprehensive logging with:
- Multiple output modes: terminal, file, or both
- Configurable log levels: debug, info, warn, error
- Structured logging with context
- Request tracing with trace IDs

## 🌐 Internationalization

Support for multiple languages:
- English (default)
- Indonesian

Translation files are located in `config/translations/`.

## 🚀 Deployment

### Docker (Recommended)

Create a `Dockerfile`:

```dockerfile
FROM golang:1.24.6-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

CMD ["./main"]
```

### Manual Deployment

1. Build the application:
```bash
go build -o xarch main.go
```

2. Run the binary:
```bash
./xarch
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Contact the development team

## 🔄 Version History

- **v1.0.0**: Initial release with HTTP and gRPC support
- Ready-to-use base code for all Go projects
- Basic health monitoring
- Multi-database support
- Internationalization support
- Clean architecture implementation

---

**Built with ❤️ by Risoftinc.**
