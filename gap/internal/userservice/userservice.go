package userservice

import (
	"context"
	"encoding/json"
	"gap/db"
	"gap/internal/ids"
)

type GuestUserSessionCreator interface {
	CreateGuestUser(context.Context, string) (db.User, error)
	CreateSession(context.Context, db.CreateSessionParams) (string, error)
	InsertEvent(context.Context, db.InsertEventParams) (db.Event, error)
}

// Create a guest user and return a session id
func CreateGuestSession(ctx context.Context, q GuestUserSessionCreator) (string, error) {
	user, err := q.CreateGuestUser(ctx, ids.Make("user"))
	if err != nil {
		return "", err
	}
	id, err := q.CreateSession(ctx, db.CreateSessionParams{SessionID: ids.Make("session"), UserID: user.UserID})
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(map[string]string{
		"UserID":   user.UserID,
		"Username": user.Username,
	})
	if err != nil {
		return "", err
	}
	_, err = q.InsertEvent(ctx, db.InsertEventParams{
		EventID:   ids.Make("evt"),
		EventType: "GuestUserCreated",
		Payload:   payload,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}
