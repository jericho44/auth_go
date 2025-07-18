# JWT Authentication API

A comprehensive Go-based authentication API with JWT tokens, built using clean architecture principles and modern Go practices.

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

### File Management

- ✅ Single and multiple file uploads
- ✅ File type validation and size limits
- ✅ Public/private file sharing
- ✅ File metadata management
- ✅ Secure file storage and retrieval

### Technical Features

- ✅ Database migrations
- ✅ **Database transactions** for data consistency
- ✅ Swagger API documentation
- ✅ CORS middleware
- ✅ Request logging
- ✅ Environment-based configuration
- ✅ Repository pattern for data access
- ✅ Dependency injection
- ✅ Clean architecture implementation
- ✅ ACID compliance with transaction support

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

#### User Management

```
GET    /api/user/profile        - Get user profile
PUT    /api/user/profile        - Update user profile
POST   /api/user/change-password - Change password
DELETE /api/user/account        - Delete account
GET    /api/user/stats          - Get user statistics
GET    /api/users               - List all users (paginated)
GET    /api/users/search        - Search users
```

#### File Management

```
POST   /api/files/upload        - Upload single file
POST   /api/files/upload/multiple - Upload multiple files
GET    /api/files               - Get user's files (paginated)
GET    /api/files/public        - Get public files
GET    /api/files/:id           - Get file details
PUT    /api/files/:id           - Update file metadata
DELETE /api/files/:id           - Delete file
GET    /api/files/download/:filename - Download file
GET    /api/files/type/:type    - Get files by type
```

## 🛠️ Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
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

5. **Start the server**

   ```bash
   go run main.go
   ```

   The application will automatically run database migrations on startup using GORM AutoMigrate.

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

### Files Table

```sql
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    original_name VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) UNIQUE NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100),
    file_type VARCHAR(50),
    user_id INTEGER NOT NULL REFERENCES users(id),
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## 🧪 Testing

The server will start on `http://localhost:8080` with the following output:

```
Connected to database successfully
Server starting on :8080
Swagger documentation available at: http://localhost:8080/swagger/
```

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

The application uses GORM AutoMigrate which automatically handles database schema updates when you start the server. This ensures your database schema stays in sync with your Go models.

### Manual Migration Commands (Optional)

If you prefer manual control over migrations:

```bash
# Run migrations
go run cmd/migrate.go -action=up

# Rollback migrations
go run cmd/migrate.go -action=down

# Check migration status
go run cmd/migrate.go -action=version
```

## 🔄 Database Transactions

This application implements comprehensive database transaction support to ensure data consistency and ACID compliance.

### Transaction Features

- **Atomicity**: All operations in a transaction succeed or fail together
- **Consistency**: Database remains in a valid state after each transaction
- **Isolation**: Concurrent transactions don't interfere with each other
- **Durability**: Committed transactions are permanently saved

### Transactional Operations

#### Repository Level

All write operations use transactions:

- User creation, updates, and deletion
- Authentication token management
- File record management
- Login attempt logging

#### Service Level

Complex operations spanning multiple database calls:

- **Password Reset**: Atomically updates password and deletes reset token
- **File Upload**: Ensures file record creation matches file system storage
- **User Registration**: Validates and creates user with proper rollback

### Example Usage

```go
// Repository level transaction (automatic)
func (ur *UserRepository) Create(user *models.User) error {
    return ur.db.Transaction(func(tx *gorm.DB) error {
        result := tx.Create(user)
        return result.Error
    })
}

// Service level transaction (complex operations)
func (as *AuthService) ResetPassword(token, newPassword string) error {
    db := as.getUserDB()
    return db.Transaction(func(tx *gorm.DB) error {
        // Update password
        if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error; err != nil {
            return err
        }

        // Delete token
        if err := tx.Where("token = ?", token).Delete(&models.PasswordResetToken{}).Error; err != nil {
            return err
        }

        return nil
    })
}
```

### Benefits

- **Data Integrity**: Prevents partial updates that could corrupt data
- **Error Recovery**: Automatic rollback on any operation failure
- **Concurrent Safety**: Proper isolation between simultaneous operations
- **Audit Trail**: Complete operation success/failure tracking

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
│   ├── file.go
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
│   ├── file.go
│   └── user.go
├── repositories/          # Data access layer
│   ├── auth_repository.go
│   ├── file_repository.go
│   └── user_repository.go
├── routes/                # Route definitions
│   ├── api.go
│   ├── auth.go
│   ├── file.go
│   ├── routes.go
│   └── user.go
├── services/              # Business services
│   ├── auth_service.go
│   ├── file_service.go
│   └── user_service.go
├── uploads/               # File upload directory
├── utils/                 # Utilities
│   ├── jwt.go
│   ├── response.go
│   └── validation.go
├── .env                   # Environment variables
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with salt
- **JWT Tokens**: Configurable expiration with HS256 signing
- **Rate Limiting**: Login attempt tracking and protection
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: GORM ORM with parameterized queries
- **CORS**: Cross-origin request handling
- **Authentication Middleware**: JWT token validation for protected routes
- **Soft Deletes**: GORM soft delete for data integrity
- **File Upload Security**: File type validation, size limits, and secure storage
- **Database Transactions**: ACID compliance for data consistency

## 🚀 Production Deployment

1. **Environment Variables**: Use secure values for production
2. **Database**: Use managed PostgreSQL service (AWS RDS, Google Cloud SQL, etc.)
3. **JWT Secret**: Generate cryptographically secure secret (32+ characters)
4. **HTTPS**: Enable TLS/SSL with reverse proxy (nginx, Cloudflare)
5. **Rate Limiting**: Configure appropriate limits for your use case
6. **Monitoring**: Add logging and monitoring (Prometheus, Grafana)
7. **Backup**: Regular database backups and disaster recovery
8. **Container**: Use Docker for consistent deployments
9. **Load Balancing**: Use load balancer for high availability

## 📝 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📚 Additional Documentation

- **[API Examples](docs/api_examples.md)** - Comprehensive API usage examples with cURL and JavaScript
- **[Adding New Features](docs/adding_new_features.md)** - Step-by-step guide for extending the application
- **[Transaction Implementation](docs/transaction_implementation.md)** - Database transaction architecture and usage
- **[File Upload Guide](docs/file_upload_guide.md)** - File management system documentation
- **[GORM Integration](docs/gorm_integration.md)** - Database ORM patterns and best practices

## 🔍 Health Check

The application includes a health check endpoint for monitoring:

```bash
curl -X GET http://localhost:8080/health
```

**Response:**

```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "database": "connected",
  "version": "1.0.0"
}
```

## 📞 Support

For questions or issues, please open a GitHub issue.
