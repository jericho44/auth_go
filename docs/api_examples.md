# API Usage Examples

This document provides comprehensive examples of how to use the JWT Authentication API with various HTTP clients.

## Authentication Flow

### 1. User Registration

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePassword123!"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. User Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePassword123!"
  }'
```

**Response:**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com"
    }
  }
}
```

## User Management

### Get User Profile

```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update User Profile

```bash
curl -X PUT http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com"
  }'
```

### Change Password

```bash
curl -X POST http://localhost:8080/api/user/change-password \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "SecurePassword123!",
    "new_password": "NewSecurePassword456!"
  }'
```

### Get User Statistics

```bash
curl -X GET http://localhost:8080/api/user/stats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user_id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "account_age_days": 30,
    "login_count": 15,
    "last_login_at": "2024-01-15T09:30:00Z",
    "created_at": "2023-12-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

## File Management

### Upload Single File

```bash
curl -X POST http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/document.pdf"
```

### Upload Multiple Files

```bash
curl -X POST http://localhost:8080/api/files/upload/multiple \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "files=@/path/to/file1.jpg" \
  -F "files=@/path/to/file2.png" \
  -F "files=@/path/to/document.pdf"
```

**Response:**

```json
{
  "success": true,
  "message": "Files uploaded successfully",
  "data": {
    "files": [
      {
        "id": 1,
        "original_name": "document.pdf",
        "file_name": "1642248600_a1b2c3d4.pdf",
        "file_size": 1024000,
        "mime_type": "application/pdf",
        "file_type": "document",
        "download_url": "/api/files/download/1642248600_a1b2c3d4.pdf",
        "is_public": false,
        "created_at": "2024-01-15T10:30:00Z"
      }
    ],
    "success_count": 1,
    "failure_count": 0,
    "failed_files": []
  }
}
```

### Get User Files

```bash
curl -X GET "http://localhost:8080/api/files?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Public Files

```bash
curl -X GET "http://localhost:8080/api/files/public?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update File Metadata

```bash
curl -X PUT http://localhost:8080/api/files/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_public": true
  }'
```

### Download File

```bash
curl -X GET http://localhost:8080/api/files/download/1642248600_a1b2c3d4.pdf \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -o downloaded_file.pdf
```

### Delete File

```bash
curl -X DELETE http://localhost:8080/api/files/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Password Reset Flow

### Request Password Reset

```bash
curl -X POST http://localhost:8080/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com"
  }'
```

### Reset Password with Token

```bash
curl -X POST http://localhost:8080/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "reset_token_from_email",
    "new_password": "NewSecurePassword789!"
  }'
```

## Admin Operations

### List All Users (Paginated)

```bash
curl -X GET "http://localhost:8080/api/users?page=1&limit=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Search Users

```bash
curl -X GET "http://localhost:8080/api/users/search?q=john&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## JavaScript/Fetch Examples

### Registration with JavaScript

```javascript
const registerUser = async (userData) => {
  try {
    const response = await fetch("http://localhost:8080/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(userData),
    });

    const result = await response.json();

    if (result.success) {
      console.log("Registration successful:", result.data);
    } else {
      console.error("Registration failed:", result.message);
    }
  } catch (error) {
    console.error("Network error:", error);
  }
};

// Usage
registerUser({
  username: "johndoe",
  email: "john@example.com",
  password: "SecurePassword123!",
});
```

### File Upload with JavaScript

```javascript
const uploadFile = async (file, token) => {
  const formData = new FormData();
  formData.append("file", file);

  try {
    const response = await fetch("http://localhost:8080/api/files/upload", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    });

    const result = await response.json();

    if (result.success) {
      console.log("Upload successful:", result.data);
    } else {
      console.error("Upload failed:", result.message);
    }
  } catch (error) {
    console.error("Upload error:", error);
  }
};
```

## Error Handling

### Common Error Responses

**Validation Error (400):**

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "Invalid email format"
    }
  ]
}
```

**Unauthorized (401):**

```json
{
  "success": false,
  "message": "Unauthorized: Invalid or expired token"
}
```

**Not Found (404):**

```json
{
  "success": false,
  "message": "Resource not found"
}
```

**Rate Limited (429):**

```json
{
  "success": false,
  "message": "Too many requests, please try again later"
}
```

**Server Error (500):**

```json
{
  "success": false,
  "message": "Internal server error"
}
```

## Best Practices

1. **Always include proper error handling** in your client applications
2. **Store JWT tokens securely** (httpOnly cookies recommended for web apps)
3. **Implement token refresh logic** before tokens expire
4. **Validate file types and sizes** on the client side before upload
5. **Use HTTPS in production** to protect sensitive data
6. **Implement proper loading states** for better user experience
7. **Handle network timeouts** and connection errors gracefully
