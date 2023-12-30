package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gap/db"
	"gap/internal/database"
	"gap/internal/ids"
	"gap/internal/server"
	"os"
	"os/exec"

	"github.com/dhowden/tag"
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
				Type:             "chat",
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
				Type:             "station",
				StationID:        payload["StationID"],
				Body:             payload["URL"] + " was requested by TODO FIXME",
				ParentID:         payload["TrackID"],
			})
			if err != nil {
				panic(err)
			}

			err := processTrackRequested(ctx, payload)
			if err != nil {
				panic(err)
			}

			err = database.CreateEvent(ctx, "TrackDownloaded", map[string]string{
				"StationID": payload["StationID"],
				"TrackID":   payload["TrackID"],
				"URL":       payload["URL"],
			})
			if err != nil {
				panic(err)
			}
		case "TrackDownloaded":
			_, err = database.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
				StationMessageID: ids.Make("sm"),
				Type:             "station",
				StationID:        payload["StationID"],
				Body:             payload["URL"] + " was added to the jukedownloaded by TODO FIXME",
				ParentID:         payload["TrackID"],
			})
			if err != nil {
				panic(err)
			}
		case "TrackStarted":
			_, err = database.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
				StationMessageID: ids.Make("sm"),
				Type:             "station",
				StationID:        payload["StationID"],
				Body:             payload["TrackID"] + " is now playing",
				ParentID:         payload["TrackID"],
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func processTrackRequested(ctx context.Context, payload map[string]string) error {
	cmd := "./yt/yt-dlp"
	stationID := payload["StationID"]
	trackID := payload["TrackID"]
	url := payload["URL"]

	// output filename pattern, does not include extension
	output := fmt.Sprintf("/tmp/%s/%s/%s", stationID, trackID, trackID)

	args := []string{
		//"--keep-video",
		"--extract-audio",
		"--audio-quality=0",
		"--audio-format=vorbis",
		"--embed-metadata",
		"--max-downloads=10",
		"--output=" + output,
		url,
	}
	out, err := exec.CommandContext(ctx, cmd, args...).Output()
	if err != nil {
		return err
	}
	fmt.Println("stdout:", string(out))

	file, err := os.Open(output + ".ogg")
	if err != nil {
		return err
	}

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return err
	}
	rawMetadata, err := json.Marshal(metadata.Raw())
	if err != nil {
		return err
	}

	_, err = database.New().Q().CreateTrack(ctx, db.CreateTrackParams{
		TrackID:     trackID,
		StationID:   stationID,
		Artist:      metadata.Artist(),
		Title:       metadata.Title(),
		RawMetadata: rawMetadata,
	})
	if err != nil {
		return err
	}

	return nil
}
