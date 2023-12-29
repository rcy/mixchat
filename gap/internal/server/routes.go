package server

import (
	_ "embed"
	"encoding/json"
	"errors"
	"gap/db"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Get("/health", s.healthHandler)

	r.Post("/create-station", s.postCreateStation)

	r.Get("/", s.stationsHandler)
	r.Get("/{slug}", s.stationHandler)
	r.Get("/{slug}/chat", s.stationHandler)
	r.Post("/{slug}/chat", s.postChatMessage)

	// liquidsoap endpoints
	r.Get("/{slug}/liq/next", s.nextHandler)

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

	// type ChatMessage struct {
	// 	CreatedAt time.Time
	// 	SentAt    string
	// 	Nick      string
	// 	Body      string
	// }
	// chatMessages := []ChatMessage{}

	//messages := chat.Fetch(station.Slug) //  s.db.Q().StationMessages(ctx, station.ID)
	messages, err := s.db.Q().StationMessages(ctx, station.StationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// for _, m := range messages {
	// 	chatMessages = append(chatMessages, ChatMessage{
	// 		CreatedAt: m.CreatedAt,
	// 		SentAt:    m.CreatedAt.Format(time.TimeOnly),
	// 		Nick:      m.Nick,
	// 		Body:      m.Body,
	// 	})
	// }
	// plays, err := s.db.Q().RecentPlays(ctx, db.RecentPlaysParams{
	// 	StationID: station.ID,
	// 	CreatedAt: pgtype.Timestamptz{Time: time.Now().Add(-time.Hour * 24), Valid: true},
	// })
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// for _, p := range plays {
	// 	chatMessages = append(chatMessages, ChatMessage{
	// 		CreatedAt: p.CreatedAt.Time,
	// 		SentAt:    p.CreatedAt.Time.Format(time.TimeOnly),
	// 		Nick:      "",
	// 		Body:      p.Artist + ": " + p.Title,
	// 	})
	// }

	// slices.SortFunc(chatMessages, func(a, b ChatMessage) int {
	// 	return b.CreatedAt.Compare(a.CreatedAt)
	// })

	// currentTrack, err := s.db.Q().CurrentTrack(ctx, station.ID)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	err = tpl.ExecuteTemplate(w, "station", struct {
		Station  db.Station
		Messages []db.StationMessage
		//		CurrentTrack db.CurrentTrackRow
	}{
		Station:  station,
		Messages: messages,
		//		CurrentTrack: currentTrack,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

type ChatEvent struct {
	ChatID    string
	Nick      string
	Body      string
	StationID string
}

func (s *Server) postChatMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event := ChatEvent{
		ChatID:    MakeID("chat"),
		StationID: station.StationID,
		Nick:      "Todo",
		Body:      r.FormValue("body"),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.db.Q().CreateEvent(r.Context(), db.CreateEventParams{
		EventID:   MakeID("evt"),
		EventType: "ChatMessageSent",
		Payload:   payload,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type StationEvent struct {
	StationID string
	Slug      string
}

func MakeID(prefix string) string {
	return prefix + "_" + ulid.Make().String()
}

func (s *Server) postCreateStation(w http.ResponseWriter, r *http.Request) {
	payload, err := json.Marshal(StationEvent{
		StationID: MakeID("stn"),
		Slug:      r.FormValue("slug"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.db.Q().CreateEvent(r.Context(), db.CreateEventParams{EventID: MakeID("evt"), EventType: "StationCreated", Payload: payload})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
