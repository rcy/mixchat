package userservice

import (
	"context"
	"encoding/json"
	"errors"
	"gap/db"
	"gap/internal/ids"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

// Create a guest user and return a session id
func CreateGuestSession(ctx context.Context, q *db.Queries) (string, error) {
	user, err := q.CreateGuestUser(ctx, ids.Make("user"))
	if err != nil {
		return "", err
	}
	id, err := q.CreateSession(ctx, ids.Make("session"), user.UserID)
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

func CreateUserSession(ctx context.Context, q *db.Queries, username string) (string, error) {
	user, err := q.UserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			user, err = q.CreateUser(ctx, ids.Make("user"), username)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	id, err := q.CreateSession(ctx, ids.Make("session"), user.UserID)
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
		EventType: "UserCreated",
		Payload:   payload,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}

const SessionCookieName = "mixchat-session"

func SetCookie(w http.ResponseWriter, sessionKey string) {
	http.SetCookie(w, &http.Cookie{
		Name:    SessionCookieName,
		Value:   sessionKey,
		Path:    "/",
		Expires: time.Now().Add(365 * 24 * time.Hour),
	})
}
