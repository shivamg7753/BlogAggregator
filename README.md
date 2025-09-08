# Blog Aggregator

A modern RSS feed aggregator built with Go, featuring real-time feed updates, user authentication, and a RESTful API.

## 🚀 Features

- **RSS Feed Management**: Add, list, and refresh RSS feeds
- **Real-time Updates**: Background job updates feeds every 5 minutes
- **User Authentication**: JWT-based authentication with secure password hashing
- **Personalized Feeds**: Users can subscribe to feeds and get personalized content
- **RESTful API**: Complete API with Swagger documentation
- **Docker Support**: Easy deployment with Docker Compose
- **PostgreSQL**: Robust database with proper relationships

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend API   │    │   Database      │
│   (React/Vue)   │◄──►│   (Go/Gin)      │◄──►│   (PostgreSQL)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │  Background     │
                       │  RSS Updater    │
                       └─────────────────┘
```

## 🛠️ Tech Stack

- **Backend**: Go 1.24, Gin, GORM
- **Database**: PostgreSQL 16
- **Authentication**: JWT with bcrypt
- **RSS Parsing**: gofeed library
- **Containerization**: Docker, Docker Compose
- **Documentation**: Swagger/OpenAPI

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Git

### 1. Clone and Setup

```bash
git clone https://github.com/shivamg7753/BlogAggregator.git
cd blogAggregator
```

### 2. Configure Environment


# Edit .env file with your settings
```

**Required Environment Variables:**
```env
POSTGRES_USER
POSTGRES_PASSWORD
POSTGRES_DB
dsn
PORT
JWT_SECRET
```

### 3. Deploy


```bash
# Start services
docker-compose up -d

# Check status
docker-compose ps
```

### 4. Access the Application

- **API**: http://localhost:8080
- **Swagger Docs**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/healthz

### 5. Frontend (React UI)

The repository includes a React frontend in `blogAggregator/frontend` built with Vite. It proxies API requests to the Go backend.

Run locally:

```bash
cd blogAggregator/frontend
npm install
npm run dev
# open the shown URL (usually http://localhost:5173)
```

Dev proxy: requests to `/api/*` are forwarded to `http://localhost:8080` (configured in `frontend/vite.config.js`).

Key UI behaviors:
- Public Posts page shows latest public posts from `/posts`.
- Login/Register available; on success, a JWT is stored and the navbar shows Logout.
- My Feed is protected; it requires login and shows your personalized feed from `/users/:id/feed`.
- Feeds page lets you create feeds, refresh them, and Subscribe (requires login). If already subscribed, the button shows “Subscribed”. If not logged in, it shows a clear error.

Environment assumptions:
- Backend runs on port 8080.
- Frontend runs on port 5173 (default Vite port).

## 📚 API Usage

### Authentication

```bash
# Register a new user
curl -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "password123"
  }'
```

### Feed Management

```bash
# Add a new feed
curl -X POST http://localhost:8080/feeds \
  -H "Content-Type: application/json" \
  -d '{
    "title": "TechCrunch",
    "url": "https://techcrunch.com/feed/"
  }'

# List all feeds
curl http://localhost:8080/feeds

# Refresh a specific feed
curl -X POST http://localhost:8080/feeds/refresh \
  -H "Content-Type: application/json" \
  -d '{"feed_id": 1}'
```

### Posts and Subscriptions

```bash
# List all posts
curl http://localhost:8080/posts

# Subscribe to a feed (requires authentication)
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "user_id": 1,
    "feed_id": 1
  }'

# Get personalized feed
curl http://localhost:8080/users/1/feed \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🔧 Development

### Local Development

```bash
# Install dependencies
go mod download

# Run locally (requires PostgreSQL)
go run cmd/main.go
```

### Database Migrations

The application automatically runs migrations on startup using GORM's AutoMigrate feature.

### Background Jobs

The RSS updater runs automatically every 5 minutes. You can also manually refresh feeds using the API.

## 🐳 Production Deployment

### Docker Compose (Recommended)

```bash
# Production deployment
docker-compose -f docker-compose.yml up -d
```

### Manual Docker

```bash
# Build image
docker build -t blog-aggregator .

# Run with PostgreSQL
docker run -d --name blog-aggregator \
  -p 8080:8080 \
  -e dsn="postgres://user:pass@host:5432/db" \
  -e JWT_SECRET="your-secret" \
  blog-aggregator
```

### Cloud Deployment

#### AWS ECS/Fargate
1. Push image to ECR
2. Create ECS task definition
3. Set up RDS PostgreSQL
4. Configure environment variables

#### Google Cloud Run
1. Build and push to GCR
2. Deploy with Cloud SQL
3. Set environment variables

#### DigitalOcean App Platform
1. Connect GitHub repository
2. Configure build settings
3. Add managed PostgreSQL database
4. Set environment variables

## 📊 Monitoring

### Health Checks

```bash
# Application health
curl http://localhost:8080/healthz

# Database connectivity
docker-compose exec app ./blogAggregator
```

### Logs

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f app
docker-compose logs -f db
```

## 🔒 Security

- **Password Hashing**: bcrypt with salt
- **JWT Tokens**: HS256 algorithm with secret key
- **Input Validation**: Gin binding validation
- **SQL Injection**: GORM ORM protection
- **CORS**: Configurable CORS settings

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check PostgreSQL is running
   - Verify connection string in .env
   - Ensure database exists

2. **JWT Token Invalid**
   - Check JWT_SECRET is set
   - Verify token format (Bearer token)
   - Check token expiration

3. **RSS Feed Not Updating**
   - Check feed URL is valid
   - Verify network connectivity
   - Check background job logs

### Support

For issues and questions:
- Check the logs: `docker-compose logs -f`
- Review API documentation: http://localhost:8080/swagger/index.html
- Open an issue on GitHub
