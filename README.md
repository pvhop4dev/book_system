# Book Management System

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/book_system)](https://goreportcard.com/report/github.com/yourusername/book_system)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A modern book management system built with Go, providing RESTful API for managing books, users, and file uploads.

## ✨ Features

- 📚 Book management (CRUD operations)
- 👥 User authentication & authorization
- 🔒 Role-based access control with Casbin
- 📁 File upload and management with MinIO
- ⚡ Redis caching for better performance
- 📊 Database migrations
- 📱 RESTful API with Swagger documentation
- 🐳 Docker & Docker Compose support

## 🚀 Tech Stack

- **Backend**: Go 1.22+
- **Database**: MySQL 8.0
- **Cache**: Redis
- **Storage**: MinIO (S3-compatible)
- **Authentication**: JWT
- **Authorization**: Casbin
- **Containerization**: Docker, Docker Compose
- **API Documentation**: Swagger

## 🛠️ Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Git

## 🏃‍♂️ Quick Start

### Using Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/book_system.git
   cd book_system
   ```

2. Copy the example config file:
   ```bash
   cp config.example.yml configs/development.yml
   ```

3. Update the configuration in `configs/development.yml` if needed.

4. Start the application with Docker Compose:
   ```bash
   docker-compose up -d
   ```

5. The application will be available at: http://localhost:8080

### Manual Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/book_system.git
   cd book_system
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the database:
   - Create a MySQL database
   - Update the database configuration in `configs/development.yml`

4. Run database migrations:
   ```bash
   go run cmd/migrate/main.go
   ```

5. Start the application:
   ```bash
   go run cmd/main.go
   ```

## 🌐 API Documentation

After starting the application, you can access the Swagger documentation at:
- Swagger UI: http://localhost:8080/swagger/index.html
- Swagger JSON: http://localhost:8080/swagger/doc.json

## 🧪 Testing

To run tests:

```bash
go test ./...
```

## 🐳 Docker Commands

- Start services: `docker-compose up -d`
- Stop services: `docker-compose down`
- View logs: `docker-compose logs -f`
- Rebuild containers: `docker-compose up -d --build`

## 📂 Project Structure

```
book_system/
├── cmd/                 # Main application entry points
├── configs/             # Configuration files
├── internal/            # Private application code
│   ├── controller/      # HTTP handlers
│   ├── domain/          # Domain models
│   ├── infrastructure/  # External services (DB, cache, storage)
│   ├── middleware/      # HTTP middlewares
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   └── transport/       # API routes and HTTP server
├── migrations/          # Database migrations
├── pkg/                 # Public packages
├── scripts/             # Utility scripts
├── .env.example         # Example environment variables
├── config.example.yml   # Example configuration
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Docker configuration
└── go.mod              # Go module definition
```

## 🔒 Environment Variables

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

Key environment variables:

- `APP_ENV`: Application environment (development, production)
- `DB_DSN`: Database connection string
- `REDIS_URL`: Redis connection URL
- `JWT_SECRET`: Secret key for JWT
- `MINIO_ENDPOINT`: MinIO server endpoint
- `MINIO_ACCESS_KEY`: MinIO access key
- `MINIO_SECRET_KEY`: MinIO secret key

## 🤝 Contributing

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 📧 Contact

Your Name - [@your_twitter](https://twitter.com/your_twitter) - your.email@example.com

Project Link: [https://github.com/yourusername/book_system](https://github.com/yourusername/book_system)

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Casbin](https://casbin.org/)
- [MinIO](https://min.io/)
- [Redis](https://redis.io/)
- [Swag](https://github.com/swaggo/swag)
