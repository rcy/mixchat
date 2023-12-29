package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gap/db"
	"gap/internal/database"
	"gap/internal/server"
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
		var payload map[string]any
		err = json.Unmarshal(event.Payload, &payload)
		if err != nil {
			panic(err)
		}
		fmt.Println("event", event.EventID, event.EventType, payload)

		switch event.EventType {
		case "StationCreated":
			_, err = database.Q().CreateStation(ctx, db.CreateStationParams{
				StationID: payload["StationID"].(string),
				Slug:      payload["Slug"].(string),
				Active:    true,
			})
			if err != nil {
				panic(err)
			}
		case "ChatMessageSent":
			_, err = database.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
				StationMessageID: server.MakeID("sm"),
				Type:             "chat",
				StationID:        payload["StationID"].(string),
				Nick:             payload["Nick"].(string),
				Body:             payload["Body"].(string),
				ParentID:         payload["ChatID"].(string),
			})
			if err != nil {
				panic(err)
			}
		}
	}
}
