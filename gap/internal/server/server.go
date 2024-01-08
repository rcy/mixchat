package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"gap/internal/database"
	"gap/internal/store"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port    int
	db      database.Service
	storage store.Store
}

func NewServer(storage store.Store) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
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
