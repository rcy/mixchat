package server

import (
	"context"
	"errors"
	"fmt"
	"gap/db"
	"gap/internal/userservice"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

// Create a guest user if a valid session cookie is not found
func (s *Server) guestUserMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		_, err := r.Cookie(sessionCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				fmt.Println("Creating guest user", err)

				tx, err := s.db.P().BeginTx(ctx, pgx.TxOptions{})
				defer tx.Rollback(ctx)

				sessionKey, err := userservice.CreateGuestSession(ctx, s.db.Q().WithTx(tx))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				err = tx.Commit(ctx)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:    sessionCookieName,
					Value:   sessionKey,
					Path:    "/",
					Expires: time.Now().Add(365 * 24 * time.Hour),
				})

				// redirect back to same path so user middleware can read the cookie
				http.Redirect(w, r, "", http.StatusSeeOther)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

type contextKey int

const userContextKey contextKey = iota

// Add current user from session to the context
func (s *Server) userMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := s.db.Q().SessionUser(ctx, cookie.Value)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting user: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, userContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func (s *Server) userFromContext(ctx context.Context) db.User {
	return ctx.Value(userContextKey).(db.User)
}

func (s *Server) requestUser(r *http.Request) db.User {
	return s.userFromContext(r.Context())
}