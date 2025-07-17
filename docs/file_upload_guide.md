# File Upload Guide

This guide explains how to use the file upload functionality in the JWT Authentication API.

## 🚀 Features

### Single File Upload

- Upload one file at a time
- Maximum file size: 10MB
- Supported formats: Images, Documents, Videos, Audio, Archives
- Automatic file type detection
- Unique filename generation
- User ownership tracking

### Multi File Upload

- Upload up to 10 files simultaneously
- Maximum total size: 50MB
- Batch processing with individual file validation
- Success/failure reporting for each file
- Partial success handling

### File Management

- View uploaded files
- Make files public/private
- Download files
- Delete files
- Filter by file type
- Pagination support

## 📁 Supported File Types

### Images

- `.jpg`, `.jpeg`, `.png`, `.gif`

### Documents

- `.pdf`, `.doc`, `.docx`, `.txt`

### Archives

- `.zip`

### Media

- `.mp4` (video), `.mp3` (audio)

## 🔧 API Endpoints

### Upload Files

#### Single File Upload

```bash
POST /api/files/upload
Content-Type: multipart/form-data
Authorization: Bearer YOUR_JWT_TOKEN

# Form data:
file: [binary file data]
```

#### Multiple Files Upload

```bash
POST /api/files/upload-multiple
Content-Type: multipart/form-data
Authorization: Bearer YOUR_JWT_TOKEN

# Form data:
files: [binary file data] (multiple files)
```

### File Management

#### Get User's Files

```bash
GET /api/files/my?page=1&limit=10
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Get File Info

```bash
GET /api/files/{id}
```

#### Update File (Make Public/Private)

```bash
PUT /api/files/{id}
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
  "is_public": true
}
```

#### Delete File

```bash
DELETE /api/files/{id}
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Download File

```bash
GET /api/files/download/{filename}
```

### Public Endpoints

#### Get Public Files

```bash
GET /api/files/public?page=1&limit=10
```

#### Get Files by Type

```bash
GET /api/files/type/{type}?page=1&limit=10
# Types: image, document, video, audio, archive, other
```

## 🧪 Testing Examples

### 1. Upload Single File with cURL

```bash
# Upload an image
curl -X POST http://localhost:8080/api/files/upload \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@/path/to/your/image.jpg"
```

### 2. Upload Multiple Files with cURL

```bash
# Upload multiple files
curl -X POST http://localhost:8080/api/files/upload-multiple \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "files=@/path/to/file1.jpg" \
  -F "files=@/path/to/file2.pdf" \
  -F "files=@/path/to/file3.txt"
```

### 3. Get Your Files

```bash
curl -X GET "http://localhost:8080/api/files/my?page=1&limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Make File Public

```bash
curl -X PUT http://localhost:8080/api/files/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_public": true}'
```

### 5. Download File

```bash
curl -X GET http://localhost:8080/api/files/download/filename.jpg \
  --output downloaded_file.jpg
```

### 6. Get Images Only

```bash
curl -X GET "http://localhost:8080/api/files/type/image?page=1&limit=10"
```

## 📝 Response Examples

### Single File Upload Response

```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "id": 1,
    "original_name": "my-photo.jpg",
    "file_name": "1642678901_a1b2c3d4.jpg",
    "file_size": 2048576,
    "mime_type": "image/jpeg",
    "file_type": "image",
    "download_url": "/api/files/download/1642678901_a1b2c3d4.jpg",
    "is_public": false,
    "created_at": "2024-01-20T10:30:00Z"
  }
}
```

### Multiple Files Upload Response

```json
{
  "success": true,
  "message": "Files uploaded successfully",
  "data": {
    "files": [
      {
        "id": 1,
        "original_name": "photo1.jpg",
        "file_name": "1642678901_a1b2c3d4.jpg",
        "file_size": 1024000,
        "mime_type": "image/jpeg",
        "file_type": "image",
        "download_url": "/api/files/download/1642678901_a1b2c3d4.jpg",
        "is_public": false,
        "created_at": "2024-01-20T10:30:00Z"
      }
    ],
    "success_count": 2,
    "failure_count": 1,
    "failed_files": ["large_video.mp4: file size exceeds maximum allowed size"]
  }
}
```

### File List Response

```json
{
  "success": true,
  "message": "Files retrieved successfully",
  "data": {
    "files": [
      {
        "id": 1,
        "original_name": "document.pdf",
        "file_name": "1642678901_a1b2c3d4.pdf",
        "file_size": 512000,
        "mime_type": "application/pdf",
        "file_type": "document",
        "download_url": "/api/files/download/1642678901_a1b2c3d4.pdf",
        "is_public": true,
        "created_at": "2024-01-20T10:30:00Z"
      }
    ],
    "total": 15,
    "page": 1,
    "limit": 10
  }
}
```

## 🔒 Security Features

### Authentication

- All upload and management endpoints require JWT authentication
- Users can only manage their own files
- Public files can be viewed by anyone

### File Validation

- File size limits (10MB per file, 50MB for multiple uploads)
- File type restrictions
- Filename sanitization
- Unique filename generation to prevent conflicts

### Access Control

- Private files by default
- Users can make files public
- Download URLs are generated dynamically
- File ownership verification

## 📂 File Storage

### Directory Structure

```
uploads/
├── 1642678901_a1b2c3d4.jpg
├── 1642678902_b2c3d4e5.pdf
└── 1642678903_c3d4e5f6.txt
```

### Database Schema

```sql
CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    original_name VARCHAR(255) NOT NULL,
    file_name VARCHAR(255) NOT NULL UNIQUE,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

## 🚨 Error Handling

### Common Errors

#### File Too Large

```json
{
  "success": false,
  "message": "File too large",
  "error": "file size exceeds maximum allowed size of 10485760 bytes"
}
```

#### Invalid File Type

```json
{
  "success": false,
  "message": "File type not allowed",
  "error": "file type .exe is not allowed"
}
```

#### No File Provided

```json
{
  "success": false,
  "message": "No file provided",
  "error": "Please select a file to upload"
}
```

#### Unauthorized Access

```json
{
  "success": false,
  "message": "You can only delete your own files",
  "error": "Forbidden"
}
```

## 🔧 Configuration

### Environment Variables

The file upload system uses the following default settings:

- **Upload Directory**: `./uploads`
- **Max File Size**: 10MB (single), 50MB (multiple)
- **Max Files**: 10 files per multi-upload
- **Allowed Extensions**: `.jpg`, `.jpeg`, `.png`, `.gif`, `.pdf`, `.doc`, `.docx`, `.txt`, `.zip`, `.mp4`, `.mp3`

### Customization

To modify these settings, update the `NewFileService` function in `services/file_service.go`:

```go
return &FileService{
    fileRepo:    fileRepo,
    userRepo:    userRepo,
    uploadPath:  "./uploads",
    maxFileSize: 10 * 1024 * 1024, // 10MB
    allowedExts: []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".txt", ".zip", ".mp4", ".mp3"},
}
```

## 🎯 Best Practices

### For Developers

1. Always validate file types and sizes
2. Use unique filenames to prevent conflicts
3. Implement proper error handling
4. Check user permissions before file operations
5. Clean up files when database operations fail

### For Users

1. Keep file sizes reasonable
2. Use descriptive filenames
3. Only make files public if necessary
4. Regularly clean up unused files
5. Be mindful of storage limits

This file upload system provides a robust, secure, and scalable solution for handling file uploads in your JWT authentication API.
