# Agnos Hospital Middleware API

Back-end candidate assignment implemented with Go, Gin, Docker Compose, Nginx, and PostgreSQL.

## Run

```bash
cp .env.example .env
docker compose up -d --build
```

API: `http://localhost:8080`

## Test

```bash
docker compose run --rm test
```

For local Go development: `go mod tidy && go test ./...`.

## Try

```bash
curl -X POST http://localhost:8080/staff/create -H 'Content-Type: application/json' \
  -d '{"username":"naruemon","password":"StrongPass123!","hospital":"hospital-a"}'

curl -X POST http://localhost:8080/staff/login -H 'Content-Type: application/json' \
  -d '{"username":"naruemon","password":"StrongPass123!","hospital":"hospital-a"}'

curl 'http://localhost:8080/patient/search?national_id=1101700203451' \
  -H 'Authorization: Bearer REPLACE_WITH_TOKEN'
```

## Documents

- [Project Structure](docs/PROJECT_STRUCTURE.md)
- [API Specification](docs/API_SPEC.md)
- [ER Diagram](docs/ER_DIAGRAM.md)
- [Combined Development Plan](docs/DEVELOPMENT_PLAN.md)

The combined development plan can be copied into the Google Doc required by the assignment.

## GitHub submission

Because the assignment is marked confidential, create a **private** GitHub repository and invite only the reviewers specified by Agnos.

```bash
git init
git add .
git commit -m "Complete Agnos back-end assignment"
git branch -M main
git remote add origin YOUR_PRIVATE_REPOSITORY_URL
git push -u origin main
```

## Security and design

- bcrypt password hashing and random opaque tokens; raw tokens are not stored.
- PostgreSQL parameterized queries.
- Hospital isolation comes from the authenticated account, never a request parameter.
- Nginx rate limiting and 1 MB body limit.
- Local Hospital A mock makes review deterministic without proprietary access.
- For production, restrict staff creation to admins, add audit logs, TLS, secret management, token revocation, and observability.
