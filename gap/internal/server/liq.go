package server

import (
	"errors"
	"fmt"
	"gap/internal/env"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var apiBase = env.MustGet("API_BASE")

func (s *Server) pullHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	track, err := s.db.Q().OldestUnplayedTrack(ctx, station.StationID)
	if errors.Is(err, pgx.ErrNoRows) {
		track, err = s.db.Q().RandomTrack(ctx, station.StationID)
		if errors.Is(err, pgx.ErrNoRows) {
			time.Sleep(time.Second)
			http.NotFound(w, r)
			return
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// I'm not sure if this should be done by the event processor, but want to avoid a race, and
	// ensure the same track doesn't get queued twice
	err = s.db.Q().IncrementTrackRotation(ctx, track.TrackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.CreateEvent(ctx, "TrackQueued", map[string]string{"StationID": station.StationID, "TrackID": track.TrackID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// key := fmt.Sprintf("%s/%s/%s.ogg", track.StationID, track.TrackID, track.TrackID)
	// uri := s.storage.URI(key)
	//url := r.URL.RawPath
	path := fmt.Sprintf("%s/%s/liq/%s", apiBase, station.Slug, track.TrackID)
	w.Write([]byte(path))
}

func (s *Server) trackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	track, err := s.db.Q().Track(ctx, chi.URLParam(r, "trackID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := os.ReadFile(fmt.Sprintf("/tmp/mixchat/%s/%s/%s.ogg", track.StationID, track.TrackID, track.TrackID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bytes)
}

func (s *Server) trackChangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filename := r.FormValue("filename")
	trackID := strings.Split(filepath.Base(filename), ".")[0]

	slug := chi.URLParam(r, "slug")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("now playing", filename, trackID)

	// Should be done by event processor?
	err = s.db.Q().IncrementTrackPlays(ctx, trackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.Q().SetStationCurrentTrack(ctx, pgtype.Text{String: trackID, Valid: true}, station.StationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.CreateEvent(ctx, "TrackStarted", map[string]string{"StationID": station.StationID, "TrackID": trackID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
