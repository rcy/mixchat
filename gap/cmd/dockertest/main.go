package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	ctx := context.TODO()

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv) //client.WithHost("tcp://localhost:2375")) // Or using TLS: "tcp://remote-host:2376"

	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	defer cli.Close()

	// // Get the list of running containers
	// containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	// if err != nil {
	// 	log.Fatalf("Error listing containers: %v", err)
	// }

	// // Print container information
	// for _, container := range containers {
	// 	fmt.Printf("Container ID: %s, Image: %s, Status: %s\n", container.ID[:10], container.Image, container.Status)
	// }

	imageName := "quickiron6536/mixchat-liquidsoap:v3"
	//imageName := "nginx:latest"

	// reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	// if err != nil {
	// 	log.Fatalf("Error pulling image: %v", err)
	// }
	// defer reader.Close()

	// fmt.Printf("Image %s pulling...\n", imageName)

	// // Wait for the image pull to complete
	// _, err = io.Copy(io.Discard, reader)
	// if err != nil {
	// 	log.Fatalf("Error reading pull response: %v", err)
	// }

	// fmt.Printf("Image %s pulled successfully\n", imageName)

	station := "rcy"

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
				fmt.Sprintf("STATION_SLUG=%s", station),
			},
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
			ExposedPorts: nat.PortSet{
				"1234/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"1234/tcp": []nat.PortBinding{
					{
						HostIP: "0.0.0.0",
						//HostPort: "1234",
					},
				},
			},
		},
		nil,
		nil,
		fmt.Sprintf("liq-%s", station),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("new container id: %v\n", response.ID)

	err = cli.ContainerStart(ctx, response.ID, container.StartOptions{})
	if err != nil {
		panic(err)
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

	// Get container logs
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true, // Follow logs (like 'tail -f')
		Timestamps: true, // Include timestamps
		Tail:       "40", // Get last 40 lines of logs
	}

	// Stream the logs
	logReader, err := cli.ContainerLogs(ctx, response.ID, options)
	if err != nil {
		log.Fatalf("Error getting container logs: %v", err)
	}
	defer logReader.Close()

	// Copy logs to stdout
	_, err = io.Copy(os.Stdout, logReader)
	if err != nil {
		log.Fatalf("Error streaming container logs: %v", err)
	}
}
