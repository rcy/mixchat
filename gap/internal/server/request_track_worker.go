package server

import (
	"context"
	"encoding/json"
	"fmt"
	"gap/db"
	"gap/internal/database"
	"gap/internal/store"
	"gap/internal/ytdlp"

	"github.com/riverqueue/river"
)

type RequestTrackArgs struct {
	Station db.Station
	User    db.User
	URL     string
	TrackID string
}

func (RequestTrackArgs) Kind() string { return "RequestTrack" }

type RequestTrackWorker struct {
	river.WorkerDefaults[RequestTrackArgs]
	Database database.Service
	Storage  store.Store
}

func (w *RequestTrackWorker) Work(ctx context.Context, job *river.Job[RequestTrackArgs]) error {
	err := download(ctx, w.Storage, downloadParams{
		StationID: job.Args.Station.StationID,
		TrackID:   job.Args.TrackID,
		URL:       job.Args.URL,
	})
	if err != nil {
		panic(err)
		// err = w.Database.CreateEvent(ctx, "TrackDownloadFailed", map[string]string{
		// 	"StationID": job.Args.Station.StationID,
		// 	"TrackID":   trackID,
		// 	"URL":       job.Args.URL,
		// 	"Error":     err.Error(),
		// })
	}

	err = w.Database.CreateEvent(ctx, "TrackDownloaded", map[string]string{
		"StationID": job.Args.Station.StationID,
		"TrackID":   job.Args.TrackID,
		"URL":       job.Args.URL,
		"Nick":      job.Args.User.Username,
	})
	if err != nil {
		panic(err)
	}

	return nil
}

type downloadParams struct {
	StationID string
	TrackID   string
	URL       string
}

func download(ctx context.Context, storage store.Store, params downloadParams) error {
	track, err := ytdlp.AudioTrackFromURL(ctx, params.URL)
	if err != nil {
		return err
	}

	rawMetadata, err := json.Marshal(track.Metadata.Raw())
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s/%s.%s", params.StationID, params.TrackID, params.TrackID, track.Format)
	if err = storage.Put(ctx, key, track.Data); err != nil {
		return err
	}

	_, err = database.New().Q().CreateTrack(ctx, db.CreateTrackParams{
		TrackID:     params.TrackID,
		StationID:   params.StationID,
		Artist:      track.Metadata.Artist(),
		Title:       track.Metadata.Title(),
		RawMetadata: rawMetadata,
	})
	if err != nil {
		return err
	}

	return nil
}
