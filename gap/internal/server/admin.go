package server

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-chi/chi/v5"
)

var (
	//go:embed admin-pages.gohtml
	adminTplContent string
	adminTpl        = template.Must(template.New("admin-pages.gohtml").Funcs(funcMap).Parse(adminTplContent))
)

type icecastStatus struct {
	URL     string
	Running bool
}

func (s *Server) adminRoute(r chi.Router) {
	r.Get("/", s.admin)
	r.Post("/start-icecast", s.startIcecast)
}

func getIcecastStatus(url string) icecastStatus {
	status := icecastStatus{URL: url}

	resp, err := http.Get(url)
	if err != nil {
		return status
	}
	if resp.StatusCode >= 300 {
		return status
	}
	status.Running = true
	return status
}

func (s *Server) admin(w http.ResponseWriter, r *http.Request) {
	status := getIcecastStatus(icecastURL)

	err := adminTpl.ExecuteTemplate(w, "home", struct {
		IcecastStatus icecastStatus
	}{
		IcecastStatus: status,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) startIcecast(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	imageName := "moul/icecast"
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
			Env:   []string{
				//"ICECAST_HOST=host.docker.internal",
				//"ICECAST_PORT=8010",
				//"ICECAST_SOURCE_PASSWORD=hackme",
			},
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
			ExposedPorts: nat.PortSet{
				"8000/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"8000/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "8010",
					},
				},
			},
		},
		nil,
		nil,
		fmt.Sprintf("mixchat-icecast-1"),
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
}
