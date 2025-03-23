package server

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
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
	r.Post("/stop-icecast", s.stopIcecast)
	r.Post("/upgrade-ytdlp", s.upgradeYtDlp)
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

func getYtDlpVersion() (string, error) {
	cmd := exec.Command("./vbin/yt-dlp", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func upgradeYtDlp() (string, error) {
	cmd := exec.Command("./vbin/yt-dlp", "--update-to", "nightly")
	output, err := cmd.Output()
	fmt.Println("output", output, "err", err)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func (s *Server) upgradeYtDlp(w http.ResponseWriter, r *http.Request) {
	output, err := upgradeYtDlp()
	if err != nil {
		http.Error(w, "upgradeYtDlp: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(output))
}

func (s *Server) admin(w http.ResponseWriter, r *http.Request) {
	status := getIcecastStatus(icecastURL)

	ytdlpVersion, err := getYtDlpVersion()
	if err != nil {
		http.Error(w, "getYtDlpVersion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = adminTpl.ExecuteTemplate(w, "home", struct {
		IcecastStatus icecastStatus
		YtDlpVersion  string
	}{
		IcecastStatus: status,
		YtDlpVersion:  ytdlpVersion,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

const (
	mixchatNetworkName   = "mixchat-radionet"
	icecastContainerName = "mixchat-icecast-1"
)

func (s *Server) startIcecast(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	summary, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	networkFound := false
	for _, ni := range summary {
		if ni.Name == mixchatNetworkName {
			networkFound = true
			break
		}
	}

	if !networkFound {
		networkResp, err := cli.NetworkCreate(ctx, mixchatNetworkName, network.CreateOptions{
			Driver: "bridge",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Network Created:", networkResp.ID)
	}

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

	// Remove existing container with the same name (if it exists)
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, c := range containers {
		if c.Names[0] == fmt.Sprintf("/%s", icecastContainerName) {
			fmt.Println("Container with this name already exists, removing it.")
			err := cli.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("Existing container removed.")
			break
		}
	}

	response, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			Cmd:   []string{},
			Env:   []string{
				//"ICECAST_HOST=host.docker.internal",
				//"ICECAST_PORT=8000",
				//"ICECAST_SOURCE_PASSWORD=hackme",
			},
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
			ExposedPorts: nat.PortSet{
				"8888/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"8888/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "8888",
					},
				},
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				mixchatNetworkName: {},
			},
		},
		nil,
		icecastContainerName,
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

func (s *Server) stopIcecast(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	err = cli.ContainerStop(ctx, icecastContainerName, container.StopOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = cli.ContainerRemove(ctx, icecastContainerName, container.RemoveOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
