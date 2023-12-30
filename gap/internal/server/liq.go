package server

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func (s *Server) pullHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	track, err := s.db.Q().OldestUnplayedTrack(ctx, station.StationID)
	if errors.Is(err, pgx.ErrNoRows) {
		track, err = s.db.Q().RandomTrack(ctx, station.StationID)
		if errors.Is(err, pgx.ErrNoRows) {
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

	filename := fmt.Sprintf("/tmp/%s/%s/%s.ogg", track.StationID, track.TrackID, track.TrackID)
	w.Write([]byte(filename))
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

	err = s.db.CreateEvent(ctx, "TrackStarted", map[string]string{"StationID": station.StationID, "TrackID": trackID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
