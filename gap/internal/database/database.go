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
}

type service struct {
	conn    *pgxpool.Pool
	queries *db.Queries
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	conn, err := pgxpool.New(context.TODO(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	queries := db.New(conn)
	s := &service{conn: conn, queries: queries}
	return s
}

func (s *service) Q() *db.Queries {
	return s.queries
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.conn.Ping(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
