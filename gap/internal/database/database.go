package database

import (
	"context"
	"fmt"
	"gap/db"
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
}

type service struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

var (
	database = os.Getenv("PGDATABASE")
	password = os.Getenv("PGPASSWORD")
	username = os.Getenv("PGUSER")
	port     = os.Getenv("PGPORT")
	host     = os.Getenv("PGHOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
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
