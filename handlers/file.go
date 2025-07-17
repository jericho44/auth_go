package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"auth-jwt/controllers"
	"auth-jwt/models"
	"auth-jwt/services"
	"auth-jwt/utils"

	"github.com/gorilla/mux"
)

var fileController *controllers.FileController
var fileService services.FileServiceInterface

func SetFileController(controller *controllers.FileController) {
	fileController = controller
}

func SetFileService(service services.FileServiceInterface) {
	fileService = service
}

// @Summary Upload a single file
// @Description Upload a single file for the authenticated user with comprehensive validation and metadata generation
// @Description
// @Description **File Requirements:**
// @Description - Maximum size: 10MB per file
// @Description - Supported formats: .jpg, .jpeg, .png, .gif, .pdf, .doc, .docx, .txt, .zip, .mp4, .mp3
// @Description - Files are private by default (can be made public later via PUT /api/files/{id})
// @Description
// @Description **Upload Process:**
// @Description 1. File validation (size, type, format)
// @Description 2. Unique filename generation (timestamp + random hash)
// @Description 3. File storage in uploads directory
// @Description 4. Database record creation with metadata
// @Description 5. Automatic file type classification
// @Description
// @Description **Response includes:**
// @Description - File ID and complete metadata
// @Description - Download URL for accessing the file
// @Description - Automatic file type classification (image, document, video, audio, archive, other)
// @Description - Original filename and generated unique filename
// @Description - File size, MIME type, and creation timestamp
// @Tags File Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload (max 10MB). Use form field name 'file'. Example: curl -F 'file=@image.jpg'"
// @Success 201 {object} controllers.FileControllerResponse "File uploaded successfully with metadata"
// @Failure 400 {object} utils.ErrorResponse "Bad request - file too large (>10MB), invalid type, or no file provided"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized - invalid or missing JWT token in Authorization header"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to save file or create database record"
// @Security BearerAuth
// @Router /api/files/upload [post]
func UploadSingleFile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		utils.BadRequest(w, "Failed to parse form", "File size too large or invalid form data")
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.BadRequest(w, "No file provided", "Please select a file to upload")
		return
	}
	defer file.Close()

	response, err := fileController.UploadSingleFile(claims.UserID, fileHeader)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "user not found"):
			utils.NotFound(w, "User not found")
		case strings.Contains(err.Error(), "file size exceeds"):
			utils.BadRequest(w, "File too large", err.Error())
		case strings.Contains(err.Error(), "file type") && strings.Contains(err.Error(), "not allowed"):
			utils.BadRequest(w, "File type not allowed", err.Error())
		case strings.Contains(err.Error(), "no file provided"):
			utils.BadRequest(w, "No file provided", "Please select a file to upload")
		default:
			utils.InternalServerError(w, "Failed to upload file")
		}
		return
	}

	utils.Created(w, response.Message, response.File)
}

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
func UploadMultipleFiles(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	// Parse multipart form
	err := r.ParseMultipartForm(50 << 20) // 50MB max for multiple files
	if err != nil {
		utils.BadRequest(w, "Failed to parse form", "Files size too large or invalid form data")
		return
	}

	// Get files from form
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		utils.BadRequest(w, "No files provided", "Please select files to upload")
		return
	}

	response, err := fileController.UploadMultipleFiles(claims.UserID, files)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "user not found"):
			utils.NotFound(w, "User not found")
		case strings.Contains(err.Error(), "maximum") && strings.Contains(err.Error(), "files allowed"):
			utils.BadRequest(w, "Too many files", err.Error())
		case strings.Contains(err.Error(), "no files provided"):
			utils.BadRequest(w, "No files provided", "Please select files to upload")
		default:
			utils.InternalServerError(w, "Failed to upload files")
		}
		return
	}

	utils.Created(w, response.Message, response.MultiFiles)
}

// @Summary Get file information by ID
// @Description Retrieve detailed information about a specific file including metadata, download URL, and file properties
// @Tags File Management
// @Produce json
// @Param id path int true "File ID" minimum(1)
// @Success 200 {object} controllers.FileControllerResponse "File information retrieved successfully"
// @Failure 400 {object} utils.ErrorResponse "Bad request - invalid file ID"
// @Failure 404 {object} utils.ErrorResponse "File not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/files/{id} [get]
func GetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid file ID", "File ID must be a number")
		return
	}

	response, err := fileController.GetFile(id)
	if err != nil {
		switch err.Error() {
		case "file not found":
			utils.NotFound(w, "File not found")
		default:
			utils.InternalServerError(w, "Failed to retrieve file")
		}
		return
	}

	utils.Success(w, "File retrieved successfully", response.File)
}

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
func GetUserFiles(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	response, err := fileController.GetUserFiles(claims.UserID, page, limit)
	if err != nil {
		utils.InternalServerError(w, "Failed to retrieve files")
		return
	}

	utils.Success(w, "Files retrieved successfully", response.Files)
}

// @Summary Get public files
// @Description Retrieve a paginated list of all files marked as public by users. No authentication required.
// @Tags File Browse
// @Produce json
// @Param page query int false "Page number for pagination" default(1) minimum(1)
// @Param limit query int false "Number of items per page" default(10) minimum(1) maximum(100)
// @Success 200 {object} controllers.FileControllerResponse "Public files retrieved successfully"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to retrieve public files"
// @Router /files/public [get]
func GetPublicFiles(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	response, err := fileController.GetPublicFiles(page, limit)
	if err != nil {
		utils.InternalServerError(w, "Failed to retrieve public files")
		return
	}

	utils.Success(w, "Public files retrieved successfully", response.Files)
}

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
func UpdateFile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid file ID", "File ID must be a number")
		return
	}

	var req models.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid request body", "Failed to parse JSON request")
		return
	}

	response, err := fileController.UpdateFile(id, claims.UserID, req)
	if err != nil {
		switch err.Error() {
		case "file not found":
			utils.NotFound(w, "File not found")
		case "unauthorized":
			utils.Forbidden(w, "You can only update your own files")
		default:
			utils.InternalServerError(w, "Failed to update file")
		}
		return
	}

	utils.Success(w, response.Message, response.File)
}

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
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequest(w, "Invalid file ID", "File ID must be a number")
		return
	}

	response, err := fileController.DeleteFile(id, claims.UserID)
	if err != nil {
		switch err.Error() {
		case "file not found":
			utils.NotFound(w, "File not found")
		case "unauthorized":
			utils.Forbidden(w, "You can only delete your own files")
		default:
			utils.InternalServerError(w, "Failed to delete file")
		}
		return
	}

	utils.Success(w, response.Message, nil)
}

// @Summary Download file
// @Description Download a file by its unique filename. The file will be served with its original name and appropriate content type headers.
// @Tags File Download
// @Produce application/octet-stream
// @Param filename path string true "Unique filename (e.g., 1642678901_a1b2c3d4.jpg)" minlength(1)
// @Success 200 {file} binary "File content with appropriate headers for download"
// @Failure 400 {object} utils.ErrorResponse "Bad request - invalid or missing filename"
// @Failure 404 {object} utils.ErrorResponse "File not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error - failed to serve file"
// @Router /files/download/{filename} [get]
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	if filename == "" {
		utils.BadRequest(w, "Invalid filename", "Filename is required")
		return
	}

	// Get file from service
	file, err := fileService.GetFileByName(filename)
	if err != nil {
		utils.NotFound(w, "File not found")
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.OriginalName+"\"")
	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(file.FileSize, 10))

	// Serve file
	http.ServeFile(w, r, file.FilePath)
}

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
// @Router /files/type/{type} [get]
func GetFilesByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileType := vars["type"]

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	response, err := fileController.GetFilesByType(fileType, page, limit)
	if err != nil {
		switch err.Error() {
		case "invalid file type":
			utils.BadRequest(w, "Invalid file type", "Valid types: image, document, video, audio, archive, other")
		default:
			utils.InternalServerError(w, "Failed to retrieve files")
		}
		return
	}

	utils.Success(w, "Files retrieved successfully", response.Files)
}
