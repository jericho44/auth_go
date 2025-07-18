# File Upload API - Swagger Documentation

This document provides detailed Swagger/OpenAPI documentation for the file upload endpoints.

## 📋 Complete Swagger Annotations

### Single File Upload

```go
// @Summary Upload a single file
// @Description Upload a single file for the authenticated user. The file will be validated for size and type, then stored securely with a unique filename.
// @Description
// @Description **File Requirements:**
// @Description - Maximum size: 10MB
// @Description - Supported formats: .jpg, .jpeg, .png, .gif, .pdf, .doc, .docx, .txt, .zip, .mp4, .mp3
// @Description - Files are private by default (can be made public later)
// @Description
// @Description **Response includes:**
// @Description - File ID and metadata
// @Description - Download URL
// @Description - File type classification
// @Tags File Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload (max 10MB)"
// @Success 201 {object} controllers.FileControllerResponse "File uploaded successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - file too large, invalid type, or no file provided"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to save file"
// @Security BearerAuth
// @Router /api/files/upload [post]
```

### Multiple Files Upload

```go
// @Summary Upload multiple files
// @Description Upload multiple files in a single request. Each file is validated individually, allowing partial success scenarios.
// @Description
// @Description **Batch Upload Limits:**
// @Description - Maximum 10 files per request
// @Description - Total size limit: 50MB
// @Description - Each file max: 10MB
// @Description - Same format restrictions as single upload
// @Description
// @Description **Response includes:**
// @Description - List of successfully uploaded files
// @Description - Success/failure counts
// @Description - Details of any failed uploads
// @Description
// @Description **Use Cases:**
// @Description - Photo gallery uploads
// @Description - Document batch processing
// @Description - Media file collections
// @Tags File Upload
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Multiple files to upload (use 'files' as field name)" multiple
// @Success 201 {object} controllers.FileControllerResponse "Files upload completed with detailed results"
// @Failure 400 {object} utils.ErrorResponse "Bad request - too many files, files too large, or no files provided"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to process files"
// @Security BearerAuth
// @Router /api/files/upload-multiple [post]
```

### File Management

```go
// @Summary Get user's uploaded files
// @Description Retrieve a paginated list of files uploaded by the authenticated user, including both public and private files
// @Tags File Management
// @Produce json
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param limit query int false "Number of items per page" default(10) minimum(1) maximum(100)
// @Success 200 {object} controllers.FileControllerResponse "User files retrieved successfully with pagination info"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to retrieve files"
// @Security BearerAuth
// @Router /api/files/my [get]
```

```go
// @Summary Update file properties
// @Description Update file properties such as making it public or private. Only the file owner can update their files.
// @Tags File Management
// @Accept json
// @Produce json
// @Param id path int true "File ID" minimum(1)
// @Param file body models.UpdateFileRequest true "File update data - currently supports changing public/private status"
// @Success 200 {object} controllers.FileControllerResponse "File updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - invalid file ID or request body"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - you can only update your own files"
// @Failure 404 {object} utils.ErrorResponse "File not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to update file"
// @Security BearerAuth
// @Router /api/files/{id} [put]
```

```go
// @Summary Delete file
// @Description Permanently delete a file owned by the authenticated user. This removes both the file record and the physical file from storage.
// @Tags File Management
// @Produce json
// @Param id path int true "File ID" minimum(1)
// @Success 200 {object} controllers.FileControllerResponse "File deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - invalid file ID"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - you can only delete your own files"
// @Failure 404 {object} utils.ErrorResponse "File not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to delete file"
// @Security BearerAuth
// @Router /api/files/{id} [delete]
```

### File Download & Browse

```go
// @Summary Download file
// @Description Download a file by its unique filename. The file will be served with its original name and appropriate content type headers.
// @Tags File Download
// @Produce application/octet-stream
// @Param filename path string true "Unique filename (e.g., 1642678901_a1b2c3d4.jpg)" minlength(1)
// @Success 200 {file} binary "File content with appropriate headers for download"
// @Failure 400 {object} utils.ErrorResponse "Bad request - invalid or missing filename"
// @Failure 404 {object} utils.ErrorResponse "File not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to serve file"
// @Router /api/files/download/{filename} [get]
```

```go
// @Summary Get public files
// @Description Retrieve a paginated list of all files marked as public by users. No authentication required.
// @Tags File Browse
// @Produce json
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param limit query int false "Number of items per page" default(10) minimum(1) maximum(100)
// @Success 200 {object} controllers.FileControllerResponse "Public files retrieved successfully"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to retrieve public files"
// @Router /api/files/public [get]
```

```go
// @Summary Get files by type
// @Description Retrieve public files filtered by file type. Useful for creating galleries or browsing specific content types.
// @Tags File Browse
// @Produce json
// @Param type path string true "File type filter" Enums(image, document, video, audio, archive, other)
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param limit query int false "Number of items per page" default(10) minimum(1) maximum(100)
// @Success 200 {object} controllers.FileControllerResponse "Files retrieved successfully by type"
// @Failure 400 {object} utils.ErrorResponse "Bad request - invalid file type (valid: image, document, video, audio, archive, other)"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to retrieve files"
// @Router /api/files/type/{type} [get]
```

## 📊 Response Schemas

### FileUploadResponse

```json
{
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
```

### MultiFileUploadResponse

```json
{
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
```

### FileListResponse

```json
{
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
```

### UpdateFileRequest

```json
{
  "is_public": true
}
```

## 🏷️ Swagger Tags

The file upload API uses the following tags for organization:

- **File Upload** - Single and multiple file upload endpoints
- **File Management** - User file management (CRUD operations)
- **File Download** - File download functionality
- **File Browse** - Public file browsing and filtering

## 🔐 Security Schemes

```yaml
securityDefinitions:
  BearerAuth:
    type: apiKey
    in: header
    name: Authorization
    description: "Enter 'Bearer' followed by a space and your JWT token"
```

## 📝 Usage Examples in Swagger UI

### Single File Upload

1. Click "Try it out"
2. Select a file using the file picker
3. Add your JWT token in the Authorization header
4. Click "Execute"

### Multiple Files Upload

1. Click "Try it out"
2. Select multiple files (hold Ctrl/Cmd)
3. Add your JWT token in the Authorization header
4. Click "Execute"

### File Management

1. Get your files using `/api/files/my`
2. Use the file ID from the response
3. Update or delete files using the respective endpoints

This comprehensive Swagger documentation provides clear, detailed information about all file upload endpoints, making it easy for developers to understand and integrate with the API.
