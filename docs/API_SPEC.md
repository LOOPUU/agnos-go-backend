# API Specification

Base URL: `http://localhost:8080`; JSON responses use `{ "data": ... }` or `{ "error": { "code": "...", "message": "..." } }`.

## POST /staff/create
Body: `{"username":"naruemon","password":"StrongPass123!","hospital":"hospital-a"}`  
Responses: `201`, `409`, `422`.

## POST /staff/login
Body: same fields as create. Returns an opaque Bearer token. Only its SHA-256 hash is stored.  
Responses: `200`, `401`, `422`.

## GET /patient/search
Header: `Authorization: Bearer <token>`  
Optional query fields: `national_id`, `passport_id`, `first_name`, `middle_name`, `last_name`, `date_of_birth`, `phone_number`, `email`. At least one is required.  
The hospital is always derived from the authenticated staff account and cannot be overridden by the caller. Maximum 100 results.  
Responses: `200`, `401`, `422`, `502`.

Example: `GET /patient/search?national_id=1101700203451`.

