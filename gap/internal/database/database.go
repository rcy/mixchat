package database

import (
	"context"
	"encoding/json"
	"fmt"
	"gap/db"
	"gap/internal/ids"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	Q() *db.Queries
	P() *pgxpool.Pool
	CreateEvent(ctx context.Context, eventType string, payloadMap map[string]string) error
}

type service struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func New() Service {
	connStr := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(context.TODO(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	queries := db.New(pool)
	s := &service{pool: pool, queries: queries}
	return s
}

func (s *service) Q() *db.Queries {
	return s.queries
}

func (s *service) P() *pgxpool.Pool {
	return s.pool
}

func (s *service) CreateEvent(ctx context.Context, eventType string, payloadMap map[string]string) error {
	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return err
	}
	_, err = s.Q().InsertEvent(ctx, db.InsertEventParams{
		EventID:   ids.Make("evt"),
		EventType: eventType,
		Payload:   payload,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.pool.Ping(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
