package controllers

import (
	"errors"
	"mime/multipart"
	"strconv"

	"auth-jwt/models"
	"auth-jwt/services"
)

type FileController struct {
	fileService services.FileServiceInterface
}

type FileControllerResponse struct {
	Message    string                          `json:"message,omitempty"`
	File       *models.FileUploadResponse      `json:"file,omitempty"`
	Files      *models.FileListResponse        `json:"files,omitempty"`
	MultiFiles *models.MultiFileUploadResponse `json:"multi_files,omitempty"`
}

func NewFileController(fileService services.FileServiceInterface) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

func (fc *FileController) UploadSingleFile(userID int, fileHeader *multipart.FileHeader) (*FileControllerResponse, error) {
	if fileHeader == nil {
		return nil, errors.New("no file provided")
	}

	file, err := fc.fileService.UploadSingleFile(uint(userID), fileHeader)
	if err != nil {
		return nil, err
	}

	return &FileControllerResponse{
		Message: "File uploaded successfully",
		File:    file,
	}, nil
}

func (fc *FileController) UploadMultipleFiles(userID int, fileHeaders []*multipart.FileHeader) (*FileControllerResponse, error) {
	if len(fileHeaders) == 0 {
		return nil, errors.New("no files provided")
	}

	if len(fileHeaders) > 10 {
		return nil, errors.New("maximum 10 files allowed per upload")
	}

	multiFiles, err := fc.fileService.UploadMultipleFiles(uint(userID), fileHeaders)
	if err != nil {
		return nil, err
	}

	message := "Files uploaded successfully"
	if multiFiles.FailureCount > 0 {
		message = "Some files failed to upload"
	}

	return &FileControllerResponse{
		Message:    message,
		MultiFiles: multiFiles,
	}, nil
}

func (fc *FileController) GetFile(id int) (*FileControllerResponse, error) {
	file, err := fc.fileService.GetFile(uint(id))
	if err != nil {
		return nil, err
	}

	return &FileControllerResponse{
		File: file,
	}, nil
}

func (fc *FileController) GetUserFiles(userID int, page, limit string) (*FileControllerResponse, error) {
	pageNum := 1
	limitNum := 10

	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitNum = l
		}
	}

	files, err := fc.fileService.GetUserFiles(uint(userID), pageNum, limitNum)
	if err != nil {
		return nil, err
	}

	return &FileControllerResponse{
		Files: files,
	}, nil
}

func (fc *FileController) GetPublicFiles(page, limit string) (*FileControllerResponse, error) {
	pageNum := 1
	limitNum := 10

	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitNum = l
		}
	}

	files, err := fc.fileService.GetPublicFiles(pageNum, limitNum)
	if err != nil {
		return nil, err
	}

	return &FileControllerResponse{
		Files: files,
	}, nil
}

func (fc *FileController) UpdateFile(id, userID int, req models.UpdateFileRequest) (*FileControllerResponse, error) {
	file, err := fc.fileService.UpdateFile(uint(id), uint(userID), req)
	if err != nil {
		return nil, err
	}

	return &FileControllerResponse{
		Message: "File updated successfully",
		File:    file,
	}, nil
}

func (fc *FileController) DeleteFile(id, userID int) (*FileControllerResponse, error) {
	err := fc.fileService.DeleteFile(uint(id), uint(userID))
	if err != nil {
		return nil, err
	}

	return &FileControllerResponse{
		Message: "File deleted successfully",
	}, nil
}

func (fc *FileController) GetFilesByType(fileType, page, limit string) (*FileControllerResponse, error) {
	// Validate file type
	validTypes := []string{models.FileTypeImage, models.FileTypeDocument, models.FileTypeVideo, models.FileTypeAudio, models.FileTypeArchive, models.FileTypeOther}
	isValid := false
	for _, validType := range validTypes {
		if fileType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, errors.New("invalid file type")
	}

	pageNum := 1
	limitNum := 10

	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p
		}
	}

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			limitNum = l
		}
	}

	files, err := fc.fileService.GetFilesByType(fileType, pageNum, limitNum)
	if err != nil {
		return nil, err
	}

	return &FileControllerResponse{
		Files: files,
	}, nil
}
