package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/LOOPUU/agnos-go-backend/internal/client"
	"github.com/LOOPUU/agnos-go-backend/internal/handler"
	"github.com/LOOPUU/agnos-go-backend/internal/repository"
	"github.com/LOOPUU/agnos-go-backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func durationEnv(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return d
}

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	repo := &repository.Repository{DB: pool}
	svc := &service.Service{
		Repo:      repo,
		Hospital:  client.New(durationEnv("HOSPITAL_API_TIMEOUT", 5*time.Second)),
		TokenTTL:  durationEnv("TOKEN_TTL", 24*time.Hour),
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal((&handler.Handler{Svc: svc}).Router().Run(":" + port))
}

