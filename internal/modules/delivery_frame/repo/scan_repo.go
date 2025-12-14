package repo

import (
	"app/internal/modules/delivery_frame/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IScanRepository interface {
	CreateScan(scan *model.Scan) error
	UpdateScanImagePath(id uuid.UUID, path string) error
	GetScanByID(id uuid.UUID) (*model.Scan, error)
	DeleteScan(id uuid.UUID) error
}

type scanRepository struct {
	db *gorm.DB
}

func NewScanRepository(db *gorm.DB) IScanRepository {
	return &scanRepository{db: db}
}

func (r *scanRepository) CreateScan(scan *model.Scan) error {
	return r.db.Create(scan).Error
}

func (r *scanRepository) UpdateScanImagePath(id uuid.UUID, path string) error {
	return r.db.Model(&model.Scan{}).Where("id = ?", id).Update("image_path", path).Error
}

func (r *scanRepository) GetScanByID(id uuid.UUID) (*model.Scan, error) {
	var scan model.Scan
	if err := r.db.First(&scan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &scan, nil
}

func (r *scanRepository) DeleteScan(id uuid.UUID) error {
	return r.db.Delete(&model.Scan{}, "id = ?", id).Error
}
