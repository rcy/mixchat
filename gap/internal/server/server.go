package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gap/internal/database"
	"gap/internal/env"
	"gap/internal/store"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port    int
	db      database.Service
	storage store.Store
}

func NewServer(storage store.Store) *http.Server {
	port, _ := strconv.Atoi(env.MustGet("PORT"))
	NewServer := &Server{
		port:    port,
		db:      database.New(),
		storage: storage,
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
