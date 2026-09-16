INSERT INTO hospitals (code, name, api_base_url) VALUES
('hospital-a', 'Hospital A', 'http://hospital-a-mock'),
('hospital-b', 'Hospital B', NULL)
ON CONFLICT (code) DO NOTHING;

INSERT INTO patients (hospital_id, patient_hn, national_id, first_name_en, last_name_en, date_of_birth, gender)
SELECT id, 'HN-B-0001', '9999999999999', 'Hospital', 'B Patient', '1985-01-01', 'F'
FROM hospitals WHERE code = 'hospital-b'
ON CONFLICT (hospital_id, patient_hn) DO NOTHING;

