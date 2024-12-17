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

	_ "github.com/joho/godotenv/autoload"
	"riverqueue.com/riverui"
)

type Server struct {
	port          int
	db            database.Service
	storage       store.Store
	riverUIServer *riverui.Server
}

func NewServer(ctx context.Context, dbService database.Service, storage store.Store, riverUIServer *riverui.Server) *http.Server {
	port, _ := strconv.Atoi(env.MustGet("PORT"))
	NewServer := &Server{
		port:          port,
		db:            dbService,
		storage:       storage,
		riverUIServer: riverUIServer,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
