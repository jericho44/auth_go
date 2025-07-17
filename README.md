# JWT Authentication API

A comprehensive Go-based authentication API with JWT tokens, built using clean architecture principles.

## 🏗️ Architecture

This project follows clean architecture with clear separation of concerns:

- **Handlers** - HTTP request/response handling
- **Controllers** - Business logic and validation
- **Services** - Application business rules
- **Repositories** - Data access layer
- **Models** - Data structures and domain entities

## 🚀 Features

### Authentication

- ✅ User registration with validation
- ✅ User login with JWT tokens
- ✅ Password reset functionality
- ✅ Rate limiting for login attempts
- ✅ Secure password hashing (bcrypt)

### User Management

- ✅ User profile management
- ✅ Password change functionality
- ✅ Account deletion
- ✅ User statistics

### Technical Features

- ✅ Database migrations
- ✅ Swagger API documentation
- ✅ CORS middleware
- ✅ Request logging
- ✅ Environment-based configuration
- ✅ Repository pattern for data access
- ✅ Dependency injection

## 📋 API Endpoints

### Public Endpoints

```
POST /auth/register          - User registration
POST /auth/login             - User login
POST /auth/forgot-password   - Password reset request
POST /auth/reset-password    - Password reset confirmation
GET  /swagger/               - API documentation
```

### Protected Endpoints (Require JWT token)

```
GET    /api/user/profile        - Get user profile
PUT    /api/user/profile        - Update user profile
POST   /api/user/change-password - Change password
DELETE /api/user/account        - Delete account
GET    /api/user/stats          - Get user statistics
```

## 🛠️ Setup

### Prerequisites

- Go 1.21+
- PostgreSQL
- Git

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd auth-go
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Setup database**

   - Create a PostgreSQL database
   - Update `.env` with your database credentials

4. **Configure environment**

   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

5. **Run migrations**

   ```bash
   go run cmd/migrate.go -action=up
   ```

6. **Start the server**
   ```bash
   go run main.go
   ```

## ⚙️ Configuration

### Environment Variables (.env)

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auth_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_TTL_HOURS=24

# Server
SERVER_PORT=8080
```

## 🗄️ Database Schema

### Users Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Password Reset Tokens Table

```sql
CREATE TABLE password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Login Attempts Table

```sql
CREATE TABLE login_attempts (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    ip_address INET,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🧪 Testing

### Manual Testing with cURL

**Register a user:**

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"password123"}'
```

**Login:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"password123"}'
```

**Access protected endpoint:**

```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Using Swagger UI

Visit `http://localhost:8080/swagger/` for interactive API testing.

## 🔧 Database Migrations

### Run migrations

```bash
go run cmd/migrate.go -action=up
```

### Rollback migrations

```bash
go run cmd/migrate.go -action=down
```

### Check migration status

```bash
go run cmd/migrate.go -action=version
```

## 📁 Project Structure

```
auth-jwt/
├── cmd/                    # Command line tools
│   └── migrate.go         # Migration tool
├── config/                # Configuration
│   └── config.go         # Environment config
├── controllers/           # Business logic controllers
│   ├── auth_controller.go
│   └── user_controller.go
├── database/              # Database connection and migrations
│   ├── db.go
│   └── migrate.go
├── docs/                  # Swagger documentation
│   ├── docs.go
│   └── swagger.json
├── handlers/              # HTTP handlers
│   ├── auth.go
│   └── user.go
├── middleware/            # HTTP middleware
│   ├── auth.go
│   ├── cors.go
│   └── logging.go
├── migrations/            # Database migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_auth_tables.up.sql
│   └── 000002_create_auth_tables.down.sql
├── models/                # Data models
│   ├── auth.go
│   └── user.go
├── repositories/          # Data access layer
│   ├── auth_repository.go
│   └── user_repository.go
├── routes/                # Route definitions
│   ├── api.go
│   ├── auth.go
│   ├── routes.go
│   └── user.go
├── services/              # Business services
│   ├── auth_service.go
│   └── user_service.go
├── utils/                 # Utilities
│   └── jwt.go
├── .env                   # Environment variables
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with salt
- **JWT Tokens**: Configurable expiration
- **Rate Limiting**: Login attempt protection
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries
- **CORS**: Cross-origin request handling

## 🚀 Production Deployment

1. **Environment Variables**: Use secure values for production
2. **Database**: Use managed PostgreSQL service
3. **JWT Secret**: Generate cryptographically secure secret
4. **HTTPS**: Enable TLS/SSL
5. **Rate Limiting**: Configure appropriate limits
6. **Monitoring**: Add logging and monitoring
7. **Backup**: Regular database backups

## 📝 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📞 Support

For questions or issues, please open a GitHub issue.
