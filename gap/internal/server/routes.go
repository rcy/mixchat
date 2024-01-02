package server

import (
	_ "embed"
	"encoding/json"
	"errors"
	"gap/db"
	"gap/internal/ids"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/rcy/durfmt"
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
	r.Get("/{slug}/audio-test-1", s.audioTest1)
	r.Get("/{slug}/audio-test-2", s.audioTest2)
	r.Post("/{slug}/requests", s.postRequest)
	r.Post("/{slug}/search", s.postSearch)
	r.Get("/{slug}/search/{searchID}", s.searchResults)

	// liquidsoap endpoints
	r.Post("/{slug}/liq/pull", s.pullHandler)
	r.Post("/{slug}/liq/trackchange", s.trackChangeHandler)
	r.Get("/{slug}/add-track", s.addTrackHandler)

	return r
}

var (
	funcMap = template.FuncMap{
		"ago": func(t time.Time) string {
			dur := time.Now().Sub(t)
			if dur < time.Minute {
				return "just now"
			}
			return durfmt.Format(dur) + " ago"
		},
	}

	//go:embed pages.gohtml
	tplContent string
	tpl        = template.Must(template.New("pages.gohtml").Funcs(funcMap).Parse(tplContent))
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

	messages, err := s.db.Q().StationMessages(ctx, station.StationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentTrack, err := s.db.Q().StationCurrentTrack(ctx, station.StationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "station", struct {
		Station      db.Station
		Messages     []db.StationMessage
		CurrentTrack db.Track
	}{
		Station:      station,
		Messages:     messages,
		CurrentTrack: currentTrack,
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

func (s *Server) postChatMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.CreateEvent(ctx, "ChatMessageSent", map[string]string{
		"StationID": station.StationID,
		"ChatID":    ids.Make("chat"),
		"Nick":      "Todo",
		"Body":      r.FormValue("body"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type CreateStationEvent struct {
	StationID string
	Slug      string
}

func (s *Server) postCreateStation(w http.ResponseWriter, r *http.Request) {
	err := s.db.CreateEvent(r.Context(), "StationCreated", map[string]string{
		"StationID": ids.Make("stn"),
		"Slug":      r.FormValue("slug"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) postRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")
	url := r.FormValue("url")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.CreateEvent(ctx, "TrackRequested", map[string]string{
		"StationID": station.StationID,
		"TrackID":   ids.Make("trk"),
		"URL":       url,
		"Nick":      "Todo",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) postSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")
	query := r.FormValue("query")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchID := ids.Make("srch")
	err = s.db.CreateEvent(ctx, "SearchSubmitted", map[string]string{
		"SearchID":  searchID,
		"StationID": station.StationID,
		"Query":     query,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.Write([]byte("searching..."))
	time.Sleep(1 * time.Second)
	http.Redirect(w, r, r.URL.Path+"/"+searchID, http.StatusSeeOther)
}

func (s *Server) searchResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	searchID := chi.URLParam(r, "searchID")
	slug := chi.URLParam(r, "slug")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	search, err := s.db.Q().Search(ctx, searchID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := s.db.Q().Results(ctx, searchID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "search-results", struct {
		Station db.Station
		Search  db.Search
		Results []db.Result
		Refresh string
	}{
		Station: station,
		Search:  search,
		Results: results,
		Refresh: r.URL.Path,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) addTrackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.ExecuteTemplate(w, "add-track", struct {
		Station db.Station
	}{
		Station: station,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) audioTest1(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "audio-test-1", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) audioTest2(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "audio-test-2", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
