CREATE TABLE IF NOT EXISTS scans (
    id UUID PRIMARY KEY,
    patient_id UUID,
    device_id UUID,
    hospital_id UUID,
    national_id VARCHAR(50),
    image_path VARCHAR(255),
    video_path VARCHAR(255),
    stone_detected BOOLEAN DEFAULT FALSE,
    hydronephrosis_detected BOOLEAN DEFAULT FALSE,
    confidence_score FLOAT,
    diagnosis_notes TEXT,
    scan_date TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
    );
