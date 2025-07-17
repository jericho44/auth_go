package repositories

import (
	"auth-jwt/models"

	"gorm.io/gorm"
)

type FileRepositoryInterface interface {
	Create(file *models.File) error
	GetByID(id uint) (*models.File, error)
	GetByFileName(fileName string) (*models.File, error)
	GetByUserID(userID uint, limit, offset int) ([]*models.File, error)
	GetPublicFiles(limit, offset int) ([]*models.File, error)
	Update(file *models.File) error
	Delete(id uint) error
	CountByUserID(userID uint) (int64, error)
	GetByFileType(fileType string, limit, offset int) ([]*models.File, error)
}

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepositoryInterface {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *FileRepository) GetByID(id uint) (*models.File, error) {
	var file models.File
	err := r.db.Preload("User").First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) GetByFileName(fileName string) (*models.File, error) {
	var file models.File
	err := r.db.Where("file_name = ?", fileName).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) GetByUserID(userID uint, limit, offset int) ([]*models.File, error) {
	var files []*models.File
	err := r.db.Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&files).Error
	return files, err
}

func (r *FileRepository) GetPublicFiles(limit, offset int) ([]*models.File, error) {
	var files []*models.File
	err := r.db.Where("is_public = ?", true).
		Preload("User").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&files).Error
	return files, err
}

func (r *FileRepository) Update(file *models.File) error {
	return r.db.Save(file).Error
}

func (r *FileRepository) Delete(id uint) error {
	return r.db.Delete(&models.File{}, id).Error
}

func (r *FileRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.File{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *FileRepository) GetByFileType(fileType string, limit, offset int) ([]*models.File, error) {
	var files []*models.File
	err := r.db.Where("file_type = ? AND is_public = ?", fileType, true).
		Preload("User").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&files).Error
	return files, err
}
