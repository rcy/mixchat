package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gap/internal/database"
	"gap/internal/env"
	"gap/internal/store"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/riverqueue/river"
	"riverqueue.com/riverui"
)

type Server struct {
	port          int
	db            database.Service
	storage       store.Store
	riverClient   *river.Client[pgx.Tx]
	riverUIServer *riverui.Server
}

var port = env.MustGet("PORT")

func NewServer(ctx context.Context, dbService database.Service, storage store.Store, riverClient *river.Client[pgx.Tx], riverUIServer *riverui.Server) *http.Server {
	portNum, _ := strconv.Atoi(port)
	server := &Server{
		port:          portNum,
		db:            dbService,
		storage:       storage,
		riverClient:   riverClient,
		riverUIServer: riverUIServer,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", server.port),
		Handler:      server.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
