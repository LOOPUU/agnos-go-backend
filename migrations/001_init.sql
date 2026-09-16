CREATE TABLE hospitals (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    api_base_url VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE staff (
    id BIGSERIAL PRIMARY KEY,
    hospital_id BIGINT NOT NULL REFERENCES hospitals(id),
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_staff_hospital_username UNIQUE (hospital_id, username)
);

CREATE TABLE access_tokens (
    id BIGSERIAL PRIMARY KEY,
    staff_id BIGINT NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    token_hash CHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE patients (
    id BIGSERIAL PRIMARY KEY,
    hospital_id BIGINT NOT NULL REFERENCES hospitals(id),
    patient_hn VARCHAR(100) NOT NULL,
    national_id VARCHAR(20),
    passport_id VARCHAR(50),
    first_name_th VARCHAR(100), middle_name_th VARCHAR(100), last_name_th VARCHAR(100),
    first_name_en VARCHAR(100), middle_name_en VARCHAR(100), last_name_en VARCHAR(100),
    date_of_birth DATE,
    phone_number VARCHAR(30),
    email VARCHAR(255),
    gender CHAR(1) CHECK (gender IN ('M','F') OR gender IS NULL),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_patient_hospital_hn UNIQUE (hospital_id, patient_hn),
    CONSTRAINT patient_identifier_required CHECK (national_id IS NOT NULL OR passport_id IS NOT NULL)
);

CREATE UNIQUE INDEX uq_patient_hospital_national ON patients(hospital_id, national_id) WHERE national_id IS NOT NULL;
CREATE UNIQUE INDEX uq_patient_hospital_passport ON patients(hospital_id, passport_id) WHERE passport_id IS NOT NULL;
CREATE INDEX idx_patient_hospital_name ON patients(hospital_id, first_name_en, last_name_en);

