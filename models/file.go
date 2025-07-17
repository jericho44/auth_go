package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	OriginalName string         `json:"original_name" gorm:"not null;size:255"`
	FileName     string         `json:"file_name" gorm:"not null;size:255;uniqueIndex"`
	FilePath     string         `json:"file_path" gorm:"not null;size:500"`
	FileSize     int64          `json:"file_size" gorm:"not null"`
	MimeType     string         `json:"mime_type" gorm:"not null;size:100"`
	FileType     string         `json:"file_type" gorm:"not null;size:50"` // image, document, video, etc.
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	User         User           `json:"user" gorm:"foreignKey:UserID"`
	IsPublic     bool           `json:"is_public" gorm:"default:false"`
	DownloadURL  string         `json:"download_url" gorm:"-"` // Not stored in DB, computed
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request/Response DTOs
type FileUploadResponse struct {
	ID           uint      `json:"id"`
	OriginalName string    `json:"original_name"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	FileType     string    `json:"file_type"`
	DownloadURL  string    `json:"download_url"`
	IsPublic     bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
}

type MultiFileUploadResponse struct {
	Files        []FileUploadResponse `json:"files"`
	SuccessCount int                  `json:"success_count"`
	FailureCount int                  `json:"failure_count"`
	FailedFiles  []string             `json:"failed_files,omitempty"`
}

type FileListResponse struct {
	Files []FileUploadResponse `json:"files"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}

type UpdateFileRequest struct {
	IsPublic *bool `json:"is_public,omitempty"`
}

// File type constants
const (
	FileTypeImage    = "image"
	FileTypeDocument = "document"
	FileTypeVideo    = "video"
	FileTypeAudio    = "audio"
	FileTypeArchive  = "archive"
	FileTypeOther    = "other"
)

// GetFileType determines file type based on MIME type
func GetFileType(mimeType string) string {
	switch {
	case len(mimeType) >= 5 && mimeType[:5] == "image":
		return FileTypeImage
	case len(mimeType) >= 5 && mimeType[:5] == "video":
		return FileTypeVideo
	case len(mimeType) >= 5 && mimeType[:5] == "audio":
		return FileTypeAudio
	case mimeType == "application/pdf" ||
		mimeType == "application/msword" ||
		mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" ||
		mimeType == "application/vnd.ms-excel" ||
		mimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" ||
		mimeType == "text/plain":
		return FileTypeDocument
	case mimeType == "application/zip" ||
		mimeType == "application/x-rar-compressed" ||
		mimeType == "application/x-7z-compressed":
		return FileTypeArchive
	default:
		return FileTypeOther
	}
}
