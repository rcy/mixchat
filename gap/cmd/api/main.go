package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gap/db"
	"gap/internal/database"
	"gap/internal/ids"
	"gap/internal/server"
	"gap/internal/ytdlp"
	"os"
	"path/filepath"
)

func main() {
	server := server.NewServer()

	ctx := context.Background()

	go process(ctx, database.New())

	fmt.Println("listening on", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

func process(ctx context.Context, database database.Service) {
	conn, err := database.P().Acquire(ctx)
	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(ctx, "listen event_inserted")
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening for events...")
	for {
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			panic(err)
		}

		if notification.Channel != "event_inserted" {
			continue
		}

		event, err := database.Q().Event(ctx, notification.Payload)
		var payload map[string]string
		err = json.Unmarshal(event.Payload, &payload)
		if err != nil {
			panic(err)
		}
		fmt.Println("event", event.EventID, event.EventType, payload)

		switch event.EventType {
		case "StationCreated":
			_, err = database.Q().CreateStation(ctx, db.CreateStationParams{
				StationID: payload["StationID"],
				Slug:      payload["Slug"],
				Active:    true,
			})
			if err != nil {
				panic(err)
			}
		case "ChatMessageSent":
			_, err = database.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
				StationMessageID: ids.Make("sm"),
				Type:             "ChatMessageSent",
				StationID:        payload["StationID"],
				Nick:             payload["Nick"],
				Body:             payload["Body"],
				ParentID:         payload["ChatID"],
			})
			if err != nil {
				panic(err)
			}
		case "TrackRequested":
			_, err = database.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
				StationMessageID: ids.Make("sm"),
				Type:             "TrackRequested",
				StationID:        payload["StationID"],
				Body:             payload["URL"],
				Nick:             payload["Nick"],
				ParentID:         payload["TrackID"],
			})
			if err != nil {
				panic(err)
			}

			err := processTrackRequested(ctx, payload)
			if err != nil {
				err = database.CreateEvent(ctx, "TrackDownloadFailed", map[string]string{
					"StationID": payload["StationID"],
					"TrackID":   payload["TrackID"],
					"URL":       payload["URL"],
					"Error":     err.Error(),
				})
				if err != nil {
					panic(err)
				}
				break
			}

			err = database.CreateEvent(ctx, "TrackDownloaded", map[string]string{
				"StationID": payload["StationID"],
				"TrackID":   payload["TrackID"],
				"URL":       payload["URL"],
				"Nick":      payload["Nick"],
			})
			if err != nil {
				panic(err)
			}
		case "TrackDownloaded":
			// Update the original TrackRequested message with a TrackDownloaded message
			// that includes some track metadata
			track, err := database.Q().Track(ctx, payload["TrackID"])
			if err != nil {
				panic(err)
			}
			m, err := database.Q().TrackRequestStationMessage(ctx, db.TrackRequestStationMessageParams{
				StationID: payload["StationID"],
				ParentID:  track.TrackID,
			})
			if err != nil {
				panic(err)
			}
			err = database.Q().UpdateStationMessage(ctx, db.UpdateStationMessageParams{
				StationMessageID: m.StationMessageID,
				Type:             "TrackDownloaded",
				Body:             fmt.Sprintf("%s, %s", track.Artist, track.Title),
			})
			if err != nil {
				panic(err)
			}
		case "TrackDownloadFailed":
			_, err = database.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
				StationMessageID: ids.Make("sm"),
				Type:             "TrackDownloadFailed",
				StationID:        payload["StationID"],
				Body:             fmt.Sprintf("Error adding %s (%s)", payload["URL"], event.EventID),
				Nick:             payload["Nick"],
				ParentID:         payload["TrackID"],
			})
			if err != nil {
				panic(err)
			}
		case "TrackStarted":
			track, err := database.Q().Track(ctx, payload["TrackID"])
			if err != nil {
				panic(err)
			}

			_, err = database.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
				StationMessageID: ids.Make("sm"),
				Type:             "TrackStarted",
				StationID:        payload["StationID"],
				Body:             fmt.Sprintf("%s, %s", track.Artist, track.Title),
				ParentID:         track.TrackID,
			})
			if err != nil {
				panic(err)
			}
		case "SearchSubmitted":
			err = database.Q().CreateSearch(ctx, db.CreateSearchParams{
				SearchID:  payload["SearchID"],
				StationID: payload["StationID"],
				Query:     payload["Query"],
			})
			if err != nil {
				panic(err)
			}

			results, err := ytdlp.Search(ctx, payload["Query"])
			if err != nil {
				err = database.CreateEvent(ctx, "SearchFailed", map[string]string{
					"StationID": payload["StationID"],
					"SearchID":  payload["SearchID"],
					"Query":     payload["Query"],
					"Error":     err.Error(),
				})
				if err != nil {
					panic(err)
				}
				break
			}

			for _, result := range results {
				err = database.Q().CreateResult(ctx, db.CreateResultParams{
					ResultID:  ids.Make("res"),
					SearchID:  payload["SearchID"],
					StationID: payload["StationID"],
					ExternID:  result.ID,
					URL:       result.WebpageURL,
					Thumbnail: result.Thumbnail,
					Title:     result.Title,
					Uploader:  result.Uploader,
					Duration:  result.Duration,
					Views:     result.ViewCount,
				})
				if err != nil {
					panic(err)
				}
			}

			err = database.Q().SetSearchStatusCompleted(ctx, payload["SearchID"])
			if err != nil {
				panic(err)
			}

			bytes, err := json.Marshal(results)
			if err != nil {
				panic(err)
			}
			err = database.CreateEvent(ctx, "SearchCompleted", map[string]string{
				"StationID": payload["StationID"],
				"SearchID":  payload["SearchID"],
				"Query":     payload["Query"],
				"Results":   string(bytes),
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func processTrackRequested(ctx context.Context, payload map[string]string) error {
	stationID := payload["StationID"]
	trackID := payload["TrackID"]
	url := payload["URL"]

	track, err := ytdlp.AudioTrackFromURL(ctx, url)
	if err != nil {
		return err
	}

	rawMetadata, err := json.Marshal(track.Metadata.Raw())
	if err != nil {
		return err
	}

	destination := fmt.Sprintf("/tmp/%s/%s/%s.%s", stationID, trackID, trackID, track.Format)
	if err = os.MkdirAll(filepath.Dir(destination), os.ModePerm); err != nil {
		return err
	}
	if err = os.WriteFile(destination, track.Data, os.ModePerm); err != nil {
		return err
	}

	_, err = database.New().Q().CreateTrack(ctx, db.CreateTrackParams{
		TrackID:     trackID,
		StationID:   stationID,
		Artist:      track.Metadata.Artist(),
		Title:       track.Metadata.Title(),
		RawMetadata: rawMetadata,
	})
	if err != nil {
		return err
	}

	return nil
}
