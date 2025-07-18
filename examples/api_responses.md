# API Response Examples

This document shows examples of the standardized API responses from the authentication system.

## Success Responses

### User Registration (201 Created)

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "message": "User registered successfully",
    "user_id": 1
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

### User Login (200 OK)

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": "2024-01-15T10:35:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

### Get User Profile (200 OK)

```json
{
  "success": true,
  "message": "User profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "timestamp": "2024-01-15T10:40:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

### User Statistics (200 OK)

```json
{
  "success": true,
  "message": "User statistics retrieved successfully",
  "data": {
    "user_id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "account_age_days": 5
  },
  "timestamp": "2024-01-15T10:45:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

### Health Check (200 OK)

```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-15T10:50:00Z",
    "version": "1.0.0",
    "uptime": "2h30m15s"
  },
  "timestamp": "2024-01-15T10:50:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

## Error Responses

### Validation Error (422 Unprocessable Entity)

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "username must be at least 3 characters long",
    "details": "Please check your input and try again"
  },
  "timestamp": "2024-01-15T10:55:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

### Unauthorized (401 Unauthorized)

```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid username or password"
  },
  "timestamp": "2024-01-15T11:00:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

### Conflict (409 Conflict)

```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "username already exists"
  },
  "timestamp": "2024-01-15T11:05:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

### Not Found (404 Not Found)

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "User profile not found"
  },
  "timestamp": "2024-01-15T11:10:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

### Bad Request (400 Bad Request)

```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid request body",
    "details": "Failed to parse JSON request"
  },
  "timestamp": "2024-01-15T11:15:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

### Internal Server Error (500 Internal Server Error)

```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_SERVER_ERROR",
    "message": "Registration failed"
  },
  "timestamp": "2024-01-15T11:20:00Z",
  "meta": {
    "version": "1.0"
  }
}
```

## Response Headers

All API responses include these headers:

```
Content-Type: application/json
X-Request-ID: a1b2c3d4e5f6g7h8
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Cache-Control: no-cache, no-store, must-revalidate
```

## Response Structure

### Success Response Structure

```json
{
  "success": true,
  "message": "Human readable success message",
  "data": {}, // Response data (can be object, array, or null)
  "meta": {
    "version": "1.0",
    "request_id": "optional",
    "page": "optional for paginated responses",
    "limit": "optional for paginated responses",
    "total": "optional for paginated responses"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Response Structure

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": "Optional additional details"
  },
  "meta": {
    "version": "1.0",
    "request_id": "optional"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```
