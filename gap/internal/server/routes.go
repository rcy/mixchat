package server

import (
	_ "embed"
	"encoding/json"
	"errors"
	"gap/db"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/exp/slices"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Get("/health", s.healthHandler)

	r.Get("/", s.stationsHandler)
	r.Get("/{slug}", s.stationHandler)
	r.Get("/{slug}/chat", s.stationHandler)
	return r
}

var (
	//go:embed pages.gohtml
	tplContent string
	tpl        = template.Must(template.New("pages.gohtml").Parse(tplContent))
)

func (s *Server) stationsHandler(w http.ResponseWriter, r *http.Request) {
	stations, err := s.db.Q().ActiveStations(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = tpl.ExecuteTemplate(w, "stations", struct{ Stations []db.Station }{Stations: stations})
}

func (s *Server) stationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	station, err := s.db.Q().Station(ctx, chi.URLParam(r, "slug"))
	if errors.Is(err, pgx.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messages, err := s.db.Q().StationMessages(ctx, station.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ChatMessage struct {
		CreatedAt time.Time
		SentAt    string
		Nick      string
		Body      string
	}
	chatMessages := []ChatMessage{}
	for _, m := range messages {
		chatMessages = append(chatMessages, ChatMessage{
			CreatedAt: m.CreatedAt.Time,
			SentAt:    m.CreatedAt.Time.Format(time.TimeOnly),
			Nick:      m.Nick.String,
			Body:      m.Body.String,
		})
	}
	plays, err := s.db.Q().RecentPlays(ctx, db.RecentPlaysParams{
		StationID: station.ID,
		CreatedAt: pgtype.Timestamptz{Time: time.Now().Add(-time.Hour * 24), Valid: true},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, p := range plays {
		chatMessages = append(chatMessages, ChatMessage{
			CreatedAt: p.CreatedAt.Time,
			SentAt:    p.CreatedAt.Time.Format(time.TimeOnly),
			Nick:      "",
			Body:      p.Artist + ": " + p.Title,
		})
	}

	slices.SortFunc(chatMessages, func(a, b ChatMessage) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	currentTrack, err := s.db.Q().CurrentTrack(ctx, station.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = tpl.ExecuteTemplate(w, "station", struct {
		Station      db.Station
		ChatMessages []ChatMessage
		CurrentTrack db.CurrentTrackRow
	}{
		Station:      station,
		ChatMessages: chatMessages,
		CurrentTrack: currentTrack,
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
