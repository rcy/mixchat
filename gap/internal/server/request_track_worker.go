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
	track, err := ytdlp.AudioTrackFromURL(ctx, job.Args.URL)
	if err != nil {
		return err
	}

	rawMetadata, err := json.Marshal(track.Metadata.Raw())
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s/%s.%s", job.Args.Station.StationID, job.Args.TrackID, job.Args.TrackID, track.Format)
	if err = w.Storage.Put(ctx, key, track.Data); err != nil {
		return err
	}

	q := w.Database.Q()

	_, err = q.CreateTrack(ctx, db.CreateTrackParams{
		TrackID:     job.Args.TrackID,
		StationID:   job.Args.Station.StationID,
		Artist:      track.Metadata.Artist(),
		Title:       track.Metadata.Title(),
		RawMetadata: rawMetadata,
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
