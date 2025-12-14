package model

import (
	"time"

	"github.com/google/uuid"
)

type Scan struct {
	ID                     uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	PatientID              uuid.UUID `json:"patient_id" gorm:"type:uuid"`
	HospitalID             uuid.UUID `json:"hospital_id" gorm:"type:uuid"`
	NationalID             string    `json:"national_id" gorm:"type:varchar(50)"`
	ImagePath              string    `json:"image_path" gorm:"type:varchar(255)"`
	VideoPath              string    `json:"video_path" gorm:"type:varchar(255)"`
	StoneDetected          bool      `json:"stone_detected" gorm:"default:false"`
	HydronephrosisDetected bool      `json:"hydronephrosis_detected" gorm:"default:false"`
	ConfidenceScore        float64   `json:"confidence_score"`
	DiagnosisNotes         string    `json:"diagnosis_notes" gorm:"type:text"`
	ScanDate               time.Time `json:"scan_date" gorm:"default:now()"`
	CreatedAt              time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"default:now()"`
}

func (Scan) TableName() string {
	return "scans"
}
