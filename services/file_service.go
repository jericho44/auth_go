package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"auth-jwt/database"
	"auth-jwt/models"
	"auth-jwt/repositories"

	"gorm.io/gorm"
)

type FileServiceInterface interface {
	UploadSingleFile(userID uint, fileHeader *multipart.FileHeader) (*models.FileUploadResponse, error)
	UploadMultipleFiles(userID uint, fileHeaders []*multipart.FileHeader) (*models.MultiFileUploadResponse, error)
	GetFile(id uint) (*models.FileUploadResponse, error)
	GetUserFiles(userID uint, page, limit int) (*models.FileListResponse, error)
	GetPublicFiles(page, limit int) (*models.FileListResponse, error)
	UpdateFile(id, userID uint, req models.UpdateFileRequest) (*models.FileUploadResponse, error)
	DeleteFile(id, userID uint) error
	GetFileByName(fileName string) (*models.File, error)
	GetFilesByType(fileType string, page, limit int) (*models.FileListResponse, error)
}

type FileService struct {
	fileRepo    repositories.FileRepositoryInterface
	userRepo    repositories.UserRepositoryInterface
	uploadPath  string
	maxFileSize int64
	allowedExts []string
}

func NewFileService(fileRepo repositories.FileRepositoryInterface, userRepo repositories.UserRepositoryInterface) FileServiceInterface {
	// Create uploads directory if it doesn't exist
	uploadPath := "./uploads"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		panic("Failed to create uploads directory: " + err.Error())
	}

	return &FileService{
		fileRepo:    fileRepo,
		userRepo:    userRepo,
		uploadPath:  uploadPath,
		maxFileSize: 10 * 1024 * 1024, // 10MB
		allowedExts: []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".txt", ".zip", ".mp4", ".mp3"},
	}
}

func (s *FileService) UploadSingleFile(userID uint, fileHeader *multipart.FileHeader) (*models.FileUploadResponse, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Validate file
	if err := s.validateFile(fileHeader); err != nil {
		return nil, err
	}

	// Generate unique filename
	fileName, err := s.generateUniqueFileName(fileHeader.Filename)
	if err != nil {
		return nil, errors.New("failed to generate filename")
	}

	// Save file to disk
	filePath := filepath.Join(s.uploadPath, fileName)
	if err := s.saveFile(fileHeader, filePath); err != nil {
		return nil, errors.New("failed to save file")
	}

	// Create file record with transaction to ensure atomicity
	file := &models.File{
		OriginalName: fileHeader.Filename,
		FileName:     fileName,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		FileType:     models.GetFileType(fileHeader.Header.Get("Content-Type")),
		UserID:       userID,
		IsPublic:     false,
	}

	// Use transaction to ensure file record is saved atomically
	db := s.getDB()
	err = db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(file).Error
	})

	if err != nil {
		// Clean up file if database save fails
		os.Remove(filePath)
		return nil, errors.New("failed to save file record")
	}

	return s.convertToResponse(file), nil
}

func (s *FileService) UploadMultipleFiles(userID uint, fileHeaders []*multipart.FileHeader) (*models.MultiFileUploadResponse, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	var uploadedFiles []models.FileUploadResponse
	var failedFiles []string
	successCount := 0
	failureCount := 0

	for _, fileHeader := range fileHeaders {
		fileResponse, err := s.UploadSingleFile(userID, fileHeader)
		if err != nil {
			failedFiles = append(failedFiles, fmt.Sprintf("%s: %s", fileHeader.Filename, err.Error()))
			failureCount++
		} else {
			uploadedFiles = append(uploadedFiles, *fileResponse)
			successCount++
		}
	}

	return &models.MultiFileUploadResponse{
		Files:        uploadedFiles,
		SuccessCount: successCount,
		FailureCount: failureCount,
		FailedFiles:  failedFiles,
	}, nil
}

func (s *FileService) GetFile(id uint) (*models.FileUploadResponse, error) {
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("file not found")
	}

	return s.convertToResponse(file), nil
}

func (s *FileService) GetUserFiles(userID uint, page, limit int) (*models.FileListResponse, error) {
	offset := (page - 1) * limit
	files, err := s.fileRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to retrieve files")
	}

	total, err := s.fileRepo.CountByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to count files")
	}

	var fileResponses []models.FileUploadResponse
	for _, file := range files {
		fileResponses = append(fileResponses, *s.convertToResponse(file))
	}

	return &models.FileListResponse{
		Files: fileResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *FileService) GetPublicFiles(page, limit int) (*models.FileListResponse, error) {
	offset := (page - 1) * limit
	files, err := s.fileRepo.GetPublicFiles(limit, offset)
	if err != nil {
		return nil, errors.New("failed to retrieve public files")
	}

	var fileResponses []models.FileUploadResponse
	for _, file := range files {
		fileResponses = append(fileResponses, *s.convertToResponse(file))
	}

	return &models.FileListResponse{
		Files: fileResponses,
		Total: int64(len(fileResponses)),
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *FileService) UpdateFile(id, userID uint, req models.UpdateFileRequest) (*models.FileUploadResponse, error) {
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("file not found")
	}

	// Check ownership
	if file.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Update fields
	if req.IsPublic != nil {
		file.IsPublic = *req.IsPublic
	}

	if err := s.fileRepo.Update(file); err != nil {
		return nil, errors.New("failed to update file")
	}

	return s.convertToResponse(file), nil
}

func (s *FileService) DeleteFile(id, userID uint) error {
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return errors.New("file not found")
	}

	// Check ownership
	if file.UserID != userID {
		return errors.New("unauthorized")
	}

	// Use transaction to ensure atomic deletion
	db := s.getDB()
	return db.Transaction(func(tx *gorm.DB) error {
		// Delete from database first
		if err := tx.Delete(&models.File{}, id).Error; err != nil {
			return err
		}

		// Delete file from disk after successful DB deletion
		if err := os.Remove(file.FilePath); err != nil {
			// Log error but don't fail the transaction since DB is already updated
			fmt.Printf("Warning: Failed to delete file from disk: %v\n", err)
		}

		return nil
	})
}

func (s *FileService) GetFileByName(fileName string) (*models.File, error) {
	return s.fileRepo.GetByFileName(fileName)
}

func (s *FileService) GetFilesByType(fileType string, page, limit int) (*models.FileListResponse, error) {
	offset := (page - 1) * limit
	files, err := s.fileRepo.GetByFileType(fileType, limit, offset)
	if err != nil {
		return nil, errors.New("failed to retrieve files by type")
	}

	var fileResponses []models.FileUploadResponse
	for _, file := range files {
		fileResponses = append(fileResponses, *s.convertToResponse(file))
	}

	return &models.FileListResponse{
		Files: fileResponses,
		Total: int64(len(fileResponses)),
		Page:  page,
		Limit: limit,
	}, nil
}

// Helper methods
func (s *FileService) validateFile(fileHeader *multipart.FileHeader) error {
	// Check file size
	if fileHeader.Size > s.maxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowed := false
	for _, allowedExt := range s.allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type %s is not allowed", ext)
	}

	return nil
}

func (s *FileService) generateUniqueFileName(originalName string) (string, error) {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Unix()

	// Generate random bytes
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Create unique filename
	fileName := fmt.Sprintf("%d_%x%s", timestamp, randomBytes, ext)
	return fileName, nil
}

func (s *FileService) saveFile(fileHeader *multipart.FileHeader, filePath string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (s *FileService) convertToResponse(file *models.File) *models.FileUploadResponse {
	return &models.FileUploadResponse{
		ID:           file.ID,
		OriginalName: file.OriginalName,
		FileName:     file.FileName,
		FileSize:     file.FileSize,
		MimeType:     file.MimeType,
		FileType:     file.FileType,
		DownloadURL:  fmt.Sprintf("/api/files/download/%s", file.FileName),
		IsPublic:     file.IsPublic,
		CreatedAt:    file.CreatedAt,
	}
}

// getDB returns the database instance for transactions
func (s *FileService) getDB() *gorm.DB {
	return database.GetDB()
}
