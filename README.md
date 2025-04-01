# Blog Platform API

[![Go Report Card](https://goreportcard.com/badge/github.com/Adejare77/blogPlatform)](https://goreportcard.com/report/github.com/Adejare77/blogPlatform)


A production-ready **Blog Platform API** built with **Go, PostgreSQL, and Redis**. Implements **RESTful best practices** with enterprise-grade **security, scalability, and performance**.

## ✨ Features

### Authentication & Authorization
- Redis-backed session management
- Middleware-based access control
- Password hashing with bcrypt

### Content Management
- Multi-status posts (**draft/published**)
- Nested comments with replies
- Polymorphic liking system
- Paginated content delivery
- Auto-generated content excerpts

### Performance
- Connection pooling (PostgreSQL/Redis)
- Index-optimized queries
- Atomic database operations
- GORM scopes for query reuse

### Operational Excellence
- Structured logging (Logrus)
- Centralized error handling
- Environment-based configuration


## 🛠 Tech Stack

### Core
- Go 1.24+
- Gin Web Framework
- GORM ORM
- Redis

### Database
- PostgreSQL 15+
- UUIDv7 primary keys
- CASCADE delete constraints


### Security
- bcrypt password hashing
- HttpOnly cookies
- Session validation middleware
- Query parameter sanitization

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker (optional)


## Configuration
### .env Template:
``` text
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=blog_admin
DB_PASSWORD=securepass123
DB_DATABASE=blogdb

# Redis
REDIS_ADDRESS=localhost:6379
REDIS_SECRETKEY=complex-secret-key
REDIS_MAX_AGE=86400

# App
CONN_MAX_LIFETIME=300s
MAX_IDLE_CONNS=20
MAX_OPEN_CONNS=100
REDIS_SIZE=10
REDIS_PASSWORD=""
```

### Installation
```bash
# Clone repository
git clone https://github.com/Adejare77/blogPlatform-gin.git
cd blogPlatform-gin

# Set up environment
cp envsample .env
nano .env  # Update with your credentials

# Install dependencies
go mod tidy

# Start server
go run cmd/cmd/main.go
```

## 📚 API Reference
## Authentication

### Register User
```http
POST /user/register
```
Request:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

### Login User
```http
POST /user/login
```
Request:

```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

### Posts
#### Get Paginated Posts
```http
GET /posts?page=2&limit=10
```
Response:
```json
{
  "data": [
    {
      "post_id": "8779ac21-4a21-41b1-9972-29c236f7602e",
      "author_name": "Jane Smith",
      "post_title": "Software Architecture Patterns",
      "content_excerpt": "Modern software architecture requires...",
      "likes": 42,
      "comments_counts": 15
    }
  ],
  "meta": {
    "page": 2,
    "limit": 10,
    "total_posts": 150,
    "_links": {
      "next": "/posts?page=3&limit=10",
      "prev": "/posts?page=1&limit=10"
    }
  }
}
```

### Comments

#### Create Comment
```http
POST /posts/{post_id}/comments
```

#### Headers:
```http
Authorization: Bearer <session_token>
Content-Type: application/json
```
### Request:
```json
{
  "content": "This is an insightful post!"
}
```

## 🗄 Database Schema
### Key Tables

#### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL
);
```

#### Posts Table

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```
#### Likes Table
```sql
CREATE TABLE likes (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    likeable_id UUID NOT NULL,
    likeable_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, likeable_id, likeable_type),
    INDEX idx_likeable (likeable_id, likeable_type)
);
```

#### Comments Table
```sql
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

## Contact

For questions or support, contact [email](rashisky007@gmail.com).