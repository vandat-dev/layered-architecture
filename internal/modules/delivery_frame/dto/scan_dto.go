package dto

import (
	"time"

	"github.com/google/uuid"
)

type ScanResponseDto struct {
	ID                     uuid.UUID `json:"id"`
	PatientID              uuid.UUID `json:"patient_id"`
	HospitalID             uuid.UUID `json:"hospital_id"`
	NationalID             string    `json:"national_id"`
	ImagePath              string    `json:"image_path"`
	VideoPath              string    `json:"video_path"`
	StoneDetected          bool      `json:"stone_detected"`
	HydronephrosisDetected bool      `json:"hydronephrosis_detected"`
	ConfidenceScore        float64   `json:"confidence_score"`
	DiagnosisNotes         string    `json:"diagnosis_notes"`
	ScanDate               time.Time `json:"scan_date"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
