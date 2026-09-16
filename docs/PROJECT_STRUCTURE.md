# Project Structure

```text
cmd/api/main.go                 composition root
internal/client/                Hospital A HTTP client
internal/handler/               Gin handlers, router, handler unit tests
internal/model/                 domain and request models
internal/repository/            PostgreSQL queries
internal/service/               business logic and authentication
migrations/                     PostgreSQL schema and seed
docker/                         Nginx configuration
mock-hospital-a/                deterministic local HIS mock
docs/                           planning documents
```

Flow: `Nginx -> Gin Handler -> Service -> Repository/PostgreSQL` and `Service -> Hospital A HIS`.

