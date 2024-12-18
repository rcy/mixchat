package server

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"gap/db"
	"gap/internal/env"
	"gap/internal/ids"
	"gap/internal/rndcolor"
	"html/template"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/rcy/durfmt"
)

var icecastURL = env.MustGet("ICECAST_URL")

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Get("/health", s.healthHandler)
	r.Get("/login", s.loginHandler)
	r.Post("/login", s.loginPostHandler)

	r.Group(func(r chi.Router) {
		r.Use(s.userMiddleware)
		// TODO: admin!
		r.Mount("/riverui", s.riverUIServer)
	})

	r.Group(func(r chi.Router) {
		//r.Use(s.guestUserMiddleware)
		r.Use(s.userMiddleware)

		r.Post("/create-station", s.postCreateStation)

		r.Get("/", s.stationsHandler)
		r.Get("/{slug}", s.stationHandler)
		r.Post("/{slug}/start-liq", s.startLiq)
		r.Get("/{slug}/now-playing", s.nowPlayingHandler)
		r.Get("/{slug}/chat", s.stationHandler)
		r.Post("/{slug}/chat", s.postChatMessage)
		r.Post("/{slug}/skip", s.postSkip)
		r.Get("/{slug}/audio-test-1", s.audioTest1)
		r.Get("/{slug}/audio-test-2", s.audioTest2)
		r.Post("/{slug}/requests", s.postRequest)
		r.Post("/{slug}/search", s.postSearch)
		r.Get("/{slug}/search/{searchID}", s.searchResults)

		r.Route("/admin", s.adminRoute)
	})

	// liquidsoap endpoints
	r.Post("/{slug}/liq/pull", s.pullHandler)
	r.Get("/{slug}/liq/{trackID}", s.trackHandler)
	r.Post("/{slug}/liq/trackchange", s.trackChangeHandler)

	return r
}

var (
	funcMap = template.FuncMap{
		"ago": func(t time.Time) string {
			dur := time.Now().Sub(t)
			if dur < time.Minute {
				return "just now"
			}
			return durfmt.Format(dur) + " ago"
		},
		"color": rndcolor.FromString,
	}

	//go:embed pages.gohtml
	tplContent string
	tpl        = template.Must(template.New("pages.gohtml").Funcs(funcMap).Parse(tplContent))
)

func (s *Server) stationsHandler(w http.ResponseWriter, r *http.Request) {
	stations, err := s.db.Q().ActiveStations(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = tpl.ExecuteTemplate(w, "stations", struct{ Stations []db.Station }{Stations: stations})
}

func (s *Server) stationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	station, err := s.db.Q().Station(ctx, chi.URLParam(r, "slug"))
	if errors.Is(err, pgx.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	messages, err := s.db.Q().StationMessages(ctx, station.StationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentTrack, err := s.db.Q().StationCurrentTrack(ctx, station.StationID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "station", struct {
		Station         db.Station
		Messages        []db.StationMessage
		CurrentTrack    db.Track
		AudioSourceURLs []string
	}{
		Station:         station,
		Messages:        messages,
		CurrentTrack:    currentTrack,
		AudioSourceURLs: []string{fmt.Sprintf("%s/%s.mp3", icecastURL, station.Slug)},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) nowPlayingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	station, err := s.db.Q().Station(ctx, chi.URLParam(r, "slug"))
	if errors.Is(err, pgx.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentTrack, err := s.db.Q().StationCurrentTrack(ctx, station.StationID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, _ := template.New("").Parse("{{.Artist}}, {{.Title}}")
	_ = t.Execute(w, currentTrack)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) postChatMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")
	user := s.requestUser(r)

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.CreateEvent(ctx, "ChatMessageSent", map[string]string{
		"StationID": station.StationID,
		"ChatID":    ids.Make("chat"),
		"UserID":    user.UserID,
		"Body":      r.FormValue("body"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.ExecuteTemplate(w, "chat-form", struct {
		Station db.Station
	}{
		Station: station,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) postSkip(w http.ResponseWriter, r *http.Request) {
	//slug := chi.URLParam(r, "slug")

	conn, err := net.DialTimeout("tcp", "localhost:1234", 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("request.dynamic.2.skip\n"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("skipx"))
}

type CreateStationEvent struct {
	StationID string
	Slug      string
}

func (s *Server) postCreateStation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := s.requestUser(r)

	stn, err := s.db.Q().CreateStation(ctx, db.CreateStationParams{
		StationID: ids.Make("stn"),
		Slug:      r.FormValue("slug"),
		UserID:    user.UserID,
		Active:    true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.CreateEvent(r.Context(), "StationCreated", map[string]string{
		"StationID": stn.StationID,
		"Slug":      stn.Slug,
		"UserID":    stn.UserID,
		"Active":    fmt.Sprint(stn.Active),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/"+stn.Slug, http.StatusSeeOther)
}

func (s *Server) postRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")
	url := r.FormValue("url")
	user := s.requestUser(r)

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trackID := ids.MakeTrackID()

	// FIXME: use transaction!

	_, err = s.db.Q().CreateStationMessage(ctx, db.CreateStationMessageParams{
		StationMessageID: ids.Make("sm"),
		Type:             "TrackRequested",
		StationID:        station.StationID,
		Body:             url,
		Nick:             user.Username,
		ParentID:         trackID,
	})
	if err != nil {
		panic(err)
	}

	_, err = s.riverClient.Insert(ctx, RequestTrackArgs{
		User:    user,
		Station: station,
		URL:     url,
		TrackID: trackID,
	}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// err = s.db.CreateEvent(ctx, "TrackRequested", map[string]string{
	// 	"StationID": station.StationID,
	// 	"TrackID":   ids.Make("trk"),
	// 	"URL":       url,
	// 	"UserID":    user.UserID,
	// })
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

func (s *Server) postSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")
	query := r.FormValue("query")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	searchID := ids.Make("srch")

	err = s.db.Q().CreateSearch(ctx, db.CreateSearchParams{
		SearchID:  searchID,
		StationID: station.StationID,
		Query:     query,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.db.CreateEvent(ctx, "SearchSubmitted", map[string]string{
		"SearchID":  searchID,
		"StationID": station.StationID,
		"Query":     query,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.URL.Path+"/"+searchID, http.StatusSeeOther)
}

func (s *Server) searchResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	searchID := chi.URLParam(r, "searchID")
	slug := chi.URLParam(r, "slug")

	station, err := s.db.Q().Station(ctx, slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	search, err := s.db.Q().Search(ctx, searchID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := s.db.Q().Results(ctx, searchID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Station   db.Station
		Search    db.Search
		Results   []db.Result
		Error     string
		Status    string
		HXGet     string
		HXTrigger string
	}{
		Station: station,
		Search:  search,
		Results: results,
	}

	switch search.Status {
	case "pending":
		data.HXGet = r.URL.Path
		data.HXTrigger = "load delay:1s"
		data.Status = fmt.Sprintf("searching for %s...", search.Query)
	case "failed":
		data.Error = "something bad happened, try again"
	}

	err = tpl.ExecuteTemplate(w, "search-results", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) audioTest1(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "audio-test-1", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) audioTest2(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "audio-test-2", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) startLiq(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	imageName := "rcy0/mixchat-liquidsoap"
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	fmt.Printf("Image %s pulling...\n", imageName)

	// Wait for the image pull to complete
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Image %s pulled successfully\n", imageName)

	response, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			Cmd:   []string{},
			Env: []string{
				"API_BASE=http://host.docker.internal:5500",
				"ICECAST_HOST=host.docker.internal",
				"ICECAST_PORT=8010",
				"ICECAST_SOURCE_PASSWORD=hackme",
				"LIQUIDSOAP_BROADCAST_PASSWORD=",
				fmt.Sprintf("STATION_SLUG=%s", slug),
			},
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
			ExposedPorts: nat.PortSet{
				"1234/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			// use random port since we have multiple containers running
			PortBindings: nat.PortMap{
				"1234/tcp": []nat.PortBinding{
					{
						HostIP: "0.0.0.0",
					},
				},
			},
		},
		nil,
		nil,
		fmt.Sprintf("liq-%s", slug),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = cli.ContainerStart(ctx, response.ID, container.StartOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	time.Sleep(5 * time.Second)

	resp, err := cli.ContainerInspect(ctx, response.ID)
	if err != nil {
		panic(err)
	}
	for containerPort, bindings := range resp.NetworkSettings.Ports {
		fmt.Println("containerPort", containerPort)
		for _, binding := range bindings {
			fmt.Printf("Container port %s is bound to host port %s\n", containerPort, binding.HostPort)
		}
	}

	fmt.Printf("Container %s started successfully\n", response.ID)
}
