package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/LOOPUU/agnos-go-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) CreateStaff(
	ctx context.Context,
	hospital, username, passwordHash string,
) (model.Staff, error) {
	var s model.Staff

	err := r.DB.QueryRow(
		ctx,
		`INSERT INTO staff(hospital_id,username,password_hash)
		SELECT id,$2,$3 FROM hospitals WHERE code=$1 RETURNING id,hospital_id,username`,
		hospital,
		username,
		passwordHash,
	).Scan(&s.ID, &s.HospitalID, &s.Username)

	return s, err
}

func (r *Repository) LoginData(
	ctx context.Context,
	hospital, username string,
) (model.Staff, string, error) {
	var s model.Staff
	var hash string

	err := r.DB.QueryRow(
		ctx,
		`SELECT s.id,s.hospital_id,h.code,s.username,h.api_base_url,s.password_hash
		FROM staff s JOIN hospitals h ON h.id=s.hospital_id WHERE h.code=$1 AND s.username=$2`,
		hospital,
		username,
	).Scan(&s.ID, &s.HospitalID, &s.HospitalCode, &s.Username, &s.APIBaseURL, &hash)

	return s, hash, err
}

func (r *Repository) SaveToken(
	ctx context.Context,
	staffID int64,
	raw string,
	expires time.Time,
) error {
	h := sha256.Sum256([]byte(raw))

	_, err := r.DB.Exec(
		ctx,
		`INSERT INTO access_tokens(staff_id,token_hash,expires_at) VALUES($1,$2,$3)`,
		staffID,
		hex.EncodeToString(h[:]),
		expires,
	)

	return err
}

func (r *Repository) Authenticate(ctx context.Context, raw string) (model.Staff, error) {
	var s model.Staff
	h := sha256.Sum256([]byte(raw))

	err := r.DB.QueryRow(
		ctx,
		`SELECT s.id,s.hospital_id,h.code,s.username,h.api_base_url FROM access_tokens t JOIN staff s ON s.id=t.staff_id JOIN hospitals h ON h.id=s.hospital_id WHERE t.token_hash=$1 AND t.expires_at>NOW()`,
		hex.EncodeToString(h[:]),
	).Scan(&s.ID, &s.HospitalID, &s.HospitalCode, &s.Username, &s.APIBaseURL)

	return s, err
}

func (r *Repository) UpsertPatient(ctx context.Context, hospitalID int64, p model.Patient) error {
	_, err := r.DB.Exec(
		ctx,
		`INSERT INTO patients(hospital_id,patient_hn,national_id,passport_id,first_name_th,middle_name_th,last_name_th,first_name_en,middle_name_en,last_name_en,date_of_birth,phone_number,email,gender)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		ON CONFLICT(hospital_id,patient_hn) DO UPDATE SET national_id=EXCLUDED.national_id,passport_id=EXCLUDED.passport_id,first_name_th=EXCLUDED.first_name_th,middle_name_th=EXCLUDED.middle_name_th,last_name_th=EXCLUDED.last_name_th,first_name_en=EXCLUDED.first_name_en,middle_name_en=EXCLUDED.middle_name_en,last_name_en=EXCLUDED.last_name_en,date_of_birth=EXCLUDED.date_of_birth,phone_number=EXCLUDED.phone_number,email=EXCLUDED.email,gender=EXCLUDED.gender,updated_at=NOW()`,
		hospitalID,
		p.PatientHN,
		p.NationalID,
		p.PassportID,
		p.FirstNameTH,
		p.MiddleNameTH,
		p.LastNameTH,
		p.FirstNameEN,
		p.MiddleNameEN,
		p.LastNameEN,
		p.DateOfBirth,
		p.PhoneNumber,
		p.Email,
		p.Gender,
	)

	return err
}

func (r *Repository) SearchPatients(
	ctx context.Context,
	hospitalID int64,
	f model.SearchFilters,
) ([]model.Patient, error) {
	where := []string{"hospital_id=$1"}
	args := []any{hospitalID}

	add := func(expr, value string) {
		if value != "" {
			args = append(args, value)
			where = append(where, fmt.Sprintf(expr, len(args)))
		}
	}

	add("national_id=$%d", f.NationalID)
	add("passport_id=$%d", f.PassportID)
	add("(first_name_th ILIKE '%%'||$%[1]d||'%%' OR first_name_en ILIKE '%%'||$%[1]d||'%%')", f.FirstName)
	add("(middle_name_th ILIKE '%%'||$%[1]d||'%%' OR middle_name_en ILIKE '%%'||$%[1]d||'%%')", f.MiddleName)
	add("(last_name_th ILIKE '%%'||$%[1]d||'%%' OR last_name_en ILIKE '%%'||$%[1]d||'%%')", f.LastName)
	add("date_of_birth=$%d", f.DateOfBirth)
	add("phone_number=$%d", f.PhoneNumber)
	add("email=$%d", f.Email)

	query := `SELECT first_name_th,middle_name_th,last_name_th,first_name_en,middle_name_en,last_name_en,date_of_birth::text,patient_hn,national_id,passport_id,phone_number,email,gender FROM patients WHERE ` + strings.Join(where, " AND ") + ` ORDER BY patient_hn LIMIT 100`

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []model.Patient{}
	for rows.Next() {
		var p model.Patient
		if err = rows.Scan(
			&p.FirstNameTH,
			&p.MiddleNameTH,
			&p.LastNameTH,
			&p.FirstNameEN,
			&p.MiddleNameEN,
			&p.LastNameEN,
			&p.DateOfBirth,
			&p.PatientHN,
			&p.NationalID,
			&p.PassportID,
			&p.PhoneNumber,
			&p.Email,
			&p.Gender,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}

	return out, rows.Err()
}
