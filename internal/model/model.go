package model

import "time"

type Staff struct {
	ID           int64   `json:"id"`
	HospitalID   int64   `json:"hospital_id"`
	HospitalCode string  `json:"hospital_code,omitempty"`
	Username     string  `json:"username"`
	APIBaseURL   *string `json:"-"`
}

type Patient struct {
	FirstNameTH  *string `json:"first_name_th"`
	MiddleNameTH *string `json:"middle_name_th"`
	LastNameTH   *string `json:"last_name_th"`
	FirstNameEN  *string `json:"first_name_en"`
	MiddleNameEN *string `json:"middle_name_en"`
	LastNameEN   *string `json:"last_name_en"`
	DateOfBirth  *string `json:"date_of_birth"`
	PatientHN    string  `json:"patient_hn"`
	NationalID   *string `json:"national_id"`
	PassportID   *string `json:"passport_id"`
	PhoneNumber  *string `json:"phone_number"`
	Email        *string `json:"email"`
	Gender       *string `json:"gender"`
}

type SearchFilters struct {
	NationalID  string `form:"national_id"`
	PassportID  string `form:"passport_id"`
	FirstName   string `form:"first_name"`
	MiddleName  string `form:"middle_name"`
	LastName    string `form:"last_name"`
	DateOfBirth string `form:"date_of_birth"`
	PhoneNumber string `form:"phone_number"`
	Email       string `form:"email"`
}

func (f SearchFilters) Empty() bool {
	return f.NationalID == "" &&
		f.PassportID == "" &&
		f.FirstName == "" &&
		f.MiddleName == "" &&
		f.LastName == "" &&
		f.DateOfBirth == "" &&
		f.PhoneNumber == "" &&
		f.Email == ""
}

type Token struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

