# Agnos Hospital Middleware API

แบบทดสอบสำหรับผู้สมัครตำแหน่ง Back-end Developer พัฒนาด้วย **Go, Gin, PostgreSQL, Docker Compose และ Nginx**

โปรเจกต์นี้เป็น API Middleware สำหรับจัดการการยืนยันตัวตนของเจ้าหน้าที่ (Staff) และค้นหาข้อมูลผู้ป่วยผ่านระบบของโรงพยาบาล

## เทคโนโลยีที่ใช้

* Go 1.25
* Gin Web Framework
* PostgreSQL 16
* Docker / Docker Compose
* Nginx
* RESTful API

## การรันระบบ

สร้างไฟล์ Environment Configuration:

```bash
cp .env.example .env
```

Build และเริ่มต้น Service ทั้งหมด:

```bash
docker compose up -d --build
```

ตรวจสอบสถานะ Service:

```bash
docker compose ps
```

API สามารถเข้าถึงได้ที่:

```text
http://localhost:8080
```

Docker Compose ประกอบด้วย:

* Nginx
* Go API
* PostgreSQL
* Hospital A Mock API

## ขั้นตอนการทำงานของ API

ขั้นตอนการทำงานหลักของระบบมีดังนี้:

1. สร้างบัญชีเจ้าหน้าที่ (Staff)
2. Login ด้วย Username และ Password ของเจ้าหน้าที่
3. ระบบส่ง Authentication Token กลับมา
4. นำ Token ไปใช้สำหรับเรียก Patient Search API
5. Middleware ตรวจสอบว่าเจ้าหน้าที่ที่ Login อยู่สังกัดโรงพยาบาลใด
6. ระบบค้นหาข้อมูลผู้ป่วยผ่านระบบของโรงพยาบาลที่เกี่ยวข้อง

สำหรับ Assignment นี้ **Hospital A ถูกจำลองเป็น Local Mock API** เพื่อให้สามารถทดสอบการทำงานของระบบได้ โดยไม่จำเป็นต้องเชื่อมต่อกับระบบจริงของโรงพยาบาล

## ทดลองเรียก API

### 1. สร้าง Staff

```bash
curl -X POST http://localhost:8080/staff/create \
  -H 'Content-Type: application/json' \
  -d '{
    "username":"naruemon",
    "password":"StrongPass123!",
    "hospital":"hospital-a"
  }'
```

### 2. Staff Login

```bash
curl -X POST http://localhost:8080/staff/login \
  -H 'Content-Type: application/json' \
  -d '{
    "username":"naruemon",
    "password":"StrongPass123!",
    "hospital":"hospital-a"
  }'
```

หลังจาก Login สำเร็จ API จะส่ง Authentication Token กลับมา

### 3. ค้นหาผู้ป่วย

นำ Token ที่ได้รับจาก Login มาแทนที่ `REPLACE_WITH_TOKEN`

```bash
curl \
  'http://localhost:8080/patient/search?national_id=1101700203451' \
  -H 'Authorization: Bearer REPLACE_WITH_TOKEN'
```

## Unit Test

ระบบมีการเขียน Unit Test สำหรับ **API Handler Layer** ครอบคลุมทั้งกรณีสำเร็จและกรณีเกิดข้อผิดพลาด

Unit Test สามารถรันผ่านได้สำเร็จโดยใช้ **Go 1.25 ผ่าน Docker**

### รัน Unit Test

```bash
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.25 \
  go test -v ./internal/handler
```

### ตรวจสอบ Test Coverage

```bash
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.25 \
  go test ./internal/handler -coverprofile=coverage.out
```

ดู Coverage แยกตาม Function:

```bash
docker run --rm \
  -v "$PWD":/app \
  -w /app \
  golang:1.25 \
  go tool cover -func=coverage.out
```

### ผลการทดสอบ

```text
Unit Test: Passed
Handler Test Coverage: 100.0%
```

ผล Coverage ของ Handler:

```text
fail          100.0%
validate      100.0%
CreateStaff   100.0%
Login         100.0%
Search        100.0%
Router        100.0%

total:        100.0%
```

Unit Test ในเวอร์ชันปัจจุบันเน้นทดสอบในส่วนของ **API Handler Layer**

## โครงสร้างโปรเจกต์

```text
agnos-go-backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── client/
│   ├── handler/
│   │   ├── handler.go
│   │   └── handler_test.go
│   ├── model/
│   ├── repository/
│   └── service/
│
├── docs/
│   ├── PROJECT_STRUCTURE.md
│   ├── API_SPEC.md
│   ├── ER_DIAGRAM.md
│   └── DEVELOPMENT_PLAN.md
│
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

### `cmd/api`

เป็นจุดเริ่มต้น (Entry Point) ของ Application ทำหน้าที่เริ่มต้น Go API Server

### `internal/handler`

ทำหน้าที่รับ HTTP Request และส่ง HTTP Response ของแต่ละ API Endpoint รวมถึงมี Unit Test สำหรับทดสอบ Handler

### `internal/service`

เก็บ Application Logic และ Business Logic ของระบบ

### `internal/repository`

ทำหน้าที่ติดต่อ PostgreSQL และจัดการการบันทึกหรือดึงข้อมูลจากฐานข้อมูล

### `internal/model`

เก็บโครงสร้างข้อมูลและ Model ที่ใช้ภายในระบบ

### `internal/client`

ทำหน้าที่ติดต่อระหว่าง Middleware กับระบบของโรงพยาบาล หรือ Mock Hospital Service

## เอกสารประกอบ

เอกสารเพิ่มเติมของโปรเจกต์อยู่ในโฟลเดอร์ `docs`:

* [Project Structure](docs/PROJECT_STRUCTURE.md)
* [API Specification](docs/API_SPEC.md)
* [ER Diagram](docs/ER_DIAGRAM.md)
* [Combined Development Plan](docs/DEVELOPMENT_PLAN.md)

Development Plan สามารถใช้เป็นเอกสารประกอบสำหรับ Candidate Assignment ได้

## ความปลอดภัยและการออกแบบระบบ

โปรเจกต์มีการออกแบบด้าน Security และโครงสร้างระบบดังนี้:

* Password ถูกเข้ารหัสแบบ Hash ด้วย bcrypt
* ไม่มีการจัดเก็บ Authentication Token แบบ Raw โดยตรง
* การ Query PostgreSQL ใช้ Parameterized Queries
* โรงพยาบาลที่ Staff สามารถเข้าถึงได้จะถูกกำหนดจากบัญชี Staff ที่ผ่านการ Authentication แล้ว แทนที่จะรับชื่อโรงพยาบาลจาก Parameter ของ Patient Search API
* Nginx มีการจำกัดอัตราการเรียก Request (Rate Limiting)
* จำกัดขนาด Request Body ไว้ที่ 1 MB
* Hospital A ถูกสร้างเป็น Local Mock Service เพื่อให้สามารถทดสอบระบบได้โดยไม่ต้องเข้าถึงระบบจริงของโรงพยาบาล

หากนำระบบไปใช้งานใน Production ควรเพิ่มมาตรการดังต่อไปนี้:

* จำกัดสิทธิ์การสร้าง Staff ให้เฉพาะ Administrator ที่ได้รับอนุญาต
* เพิ่ม Audit Log
* ใช้ HTTPS/TLS
* จัดการ Secret อย่างปลอดภัย
* เพิ่มระบบ Token Expiration และ Token Revocation
* เพิ่ม Monitoring และ Observability
* เพิ่ม Unit Test และ Integration Test ให้ครอบคลุมส่วนอื่น ๆ ของระบบมากขึ้น

## การส่ง Assignment

Repository นี้ประกอบด้วย Source Code, Docker Configuration, Unit Test, Mock Hospital Integration และเอกสารประกอบที่ใช้สำหรับ Back-end Candidate Assignment

เนื่องจาก Assignment ถูกระบุว่าเป็น **ข้อมูล Confidential** จึงควรตั้ง Repository เป็น **Private** และอนุญาตให้เข้าถึงเฉพาะผู้ตรวจที่ Agnos ระบุเท่านั้น
