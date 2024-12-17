package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gap/db"
	"gap/internal/database"
	"gap/internal/ids"
	"gap/internal/server"
	"gap/internal/store/files"
	"gap/internal/ytdlp"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/tern/v2/migrate"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"riverqueue.com/riverui"
)

func main() {
	ctx := context.Background()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("main caught signal", sig)

		os.Exit(1)
	}()

	storage := files.MustInit("/tmp/mixchat")
	// storage := space.MustInit(space.InitParams{
	// 	S3Key:       os.Getenv("S3_ACCESS_KEY"),
	// 	S3Secret:    os.Getenv("S3_SECRET_KEY"),
	// 	Endpoint:    os.Getenv("S3_ENDPOINT"),
	// 	URIEndpoint: os.Getenv("S3_URI_ENDPOINT"),
	// 	Bucket:      os.Getenv("S3_BUCKET"),
	// })

	dbService := database.New()

	poolConn, err := dbService.P().Acquire(ctx)
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewMigrator(ctx, poolConn.Conn(), "public.schema_version")
	if err != nil {
		panic(err)
	}

	migrationRoot, _ := fs.Sub(db.MigrationFiles, "migrations")
	err = m.LoadMigrations(migrationRoot)
	if err != nil {
		panic(err)
	}

	version, err := m.GetCurrentVersion(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("schema_version: %d/%d\n", version, len(m.Migrations))

	err = m.Migrate(ctx)
	if err != nil {
		panic(err)
	}

	workers := river.NewWorkers()
	river.AddWorker(workers, &DummyWorker{})
	river.AddWorker(workers, &server.RequestTrackWorker{Storage: storage, Database: dbService})

	riverClient, err := river.NewClient(riverpgxv5.New(dbService.P()), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		Workers: workers,
	})
	if err != nil {
		panic(err)
	}

	// Run the client inline. All executed jobs will inherit from ctx:
	if err := riverClient.Start(ctx); err != nil {
		panic(err)
	}

	go process(ctx, dbService)

	riverUIServer, err := riverui.NewServer(&riverui.ServerOpts{
		Client: riverClient,
		DB:     dbService.P(),
		Logger: slog.Default(),
		Prefix: "/riverui", // mount the UI and its APIs under /riverui
		// ...
	})
	if err != nil {
		panic(err)
	}
	// Start the server to initialize background processes for caching and periodic queries:
	err = riverUIServer.Start(ctx)
	if err != nil {
		panic(err)
	}

	webServer := server.NewServer(ctx, dbService, storage, riverClient, riverUIServer)

	fmt.Println("listening on", webServer.Addr)
	err = webServer.ListenAndServe()
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

		go func() {
			event, err := database.Q().Event(ctx, notification.Payload)
			var payload map[string]string
			err = json.Unmarshal(event.Payload, &payload)
			if err != nil {
				panic(err)
			}
			fmt.Println(event.EventType, event.EventID, payload)

			switch event.EventType {
			case "ChatMessageSent":
				user, err := database.Q().User(ctx, payload["UserID"])
				if err != nil {
					panic(err)
				}

				_, err = database.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
					StationMessageID: ids.Make("sm"),
					Type:             "ChatMessageSent",
					StationID:        payload["StationID"],
					Nick:             user.Username,
					Body:             payload["Body"],
					ParentID:         payload["ChatID"],
				})
				if err != nil {
					panic(err)
				}
			// case "TrackRequested":
			// 	user, err := database.Q().User(ctx, payload["UserID"])
			// 	if err != nil {
			// 		panic(err)
			// 	}

			// 	_, err = database.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
			// 		StationMessageID: ids.Make("sm"),
			// 		Type:             "TrackRequested",
			// 		StationID:        payload["StationID"],
			// 		Body:             payload["URL"],
			// 		Nick:             user.Username,
			// 		ParentID:         payload["TrackID"],
			// 	})
			// 	if err != nil {
			// 		panic(err)
			// 	}

			// 	err = processTrackRequested(ctx, str, payload)
			// 	if err != nil {
			// 		err = database.CreateEvent(ctx, "TrackDownloadFailed", map[string]string{
			// 			"StationID": payload["StationID"],
			// 			"TrackID":   payload["TrackID"],
			// 			"URL":       payload["URL"],
			// 			"Error":     err.Error(),
			// 		})
			// 		if err != nil {
			// 			panic(err)
			// 		}
			// 		break
			// 	}

			// 	err = database.CreateEvent(ctx, "TrackDownloaded", map[string]string{
			// 		"StationID": payload["StationID"],
			// 		"TrackID":   payload["TrackID"],
			// 		"URL":       payload["URL"],
			// 		"Nick":      payload["Nick"],
			// 	})
			// 	if err != nil {
			// 		panic(err)
			// 	}
			case "TrackDownloaded":
				// Update the original TrackRequested message with a TrackDownloaded message
				// that includes some track metadata
				track, err := database.Q().Track(ctx, payload["TrackID"])
				if err != nil {
					panic(err)
				}
				m, err := database.Q().TrackRequestStationMessage(ctx, payload["StationID"], track.TrackID)
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
				// Update the original TrackRequested message with a TrackDownloadFailed message
				m, err := database.Q().TrackRequestStationMessage(ctx, payload["StationID"], payload["TrackID"])
				if err != nil {
					panic(err)
				}
				err = database.Q().UpdateStationMessage(ctx, db.UpdateStationMessageParams{
					StationMessageID: m.StationMessageID,
					Type:             "TrackDownloadFailed",
					Body:             fmt.Sprintf("Error adding track (%s)", event.EventID),
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
				results, err := ytdlp.Search(ctx, payload["Query"])
				if err != nil {
					origErr := err
					err = database.Q().SetSearchStatusFailed(ctx, payload["SearchID"])
					if err != nil {
						panic(err)
					}

					err = database.CreateEvent(ctx, "SearchFailed", map[string]string{
						"StationID": payload["StationID"],
						"SearchID":  payload["SearchID"],
						"Query":     payload["Query"],
						"Error":     origErr.Error(),
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
		}()
	}
}

// func processTrackRequested(ctx context.Context, storage store.Store, payload map[string]string) error {
// 	stationID := payload["StationID"]
// 	trackID := payload["TrackID"]
// 	url := payload["URL"]

// 	track, err := ytdlp.AudioTrackFromURL(ctx, url)
// 	if err != nil {
// 		return err
// 	}

// 	rawMetadata, err := json.Marshal(track.Metadata.Raw())
// 	if err != nil {
// 		return err
// 	}

// 	key := fmt.Sprintf("%s/%s/%s.%s", stationID, trackID, trackID, track.Format)
// 	if err = storage.Put(ctx, key, track.Data); err != nil {
// 		return err
// 	}

// 	_, err = database.New().Q().CreateTrack(ctx, db.CreateTrackParams{
// 		TrackID:     trackID,
// 		StationID:   stationID,
// 		Artist:      track.Metadata.Artist(),
// 		Title:       track.Metadata.Title(),
// 		RawMetadata: rawMetadata,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
