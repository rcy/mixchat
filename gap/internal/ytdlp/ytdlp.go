package ytdlp

import (
	"context"
	"fmt"
	"gap/internal/env"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/dhowden/tag"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type format string

const (
	OGG        format = "ogg"
	executable        = "./vbin/yt-dlp"
)

type AudioTrack struct {
	Format   format
	Data     []byte
	Metadata tag.Metadata
}

func AudioTrackFromURL(ctx context.Context, url string) (*AudioTrack, error) {
	dir, err := os.MkdirTemp("", "mixchat-")
	if err != nil {
		return nil, err
	}
	if os.Getenv("KEEP_TEMP") != "true" {
		defer os.RemoveAll(dir)
	}

	fmt.Printf("ripping %s into %s\n", url, dir)

	destination := fmt.Sprintf("%s/track", dir)

	args := []string{
		"--cookies=./vbin/cookies.txt",
		"--keep-video",
		"--extract-audio",
		//"--audio-quality=0",
		"--audio-format=vorbis",
		"--embed-metadata",
		"--max-downloads=10",
		"--output=" + destination,
		url,
	}

	fmt.Printf("args=%v\n", args)

	output, err := exec.CommandContext(ctx, executable, args...).CombinedOutput()
	fmt.Println("output:", string(output))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, string(output))
	}

	file, err := os.Open(destination + ".ogg")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &AudioTrack{
		Format:   OGG,
		Data:     data,
		Metadata: metadata,
	}, nil
}

type Result struct {
	ID         string
	Thumbnail  string
	Title      string
	Uploader   string
	WebpageURL string `json:"webpage_url"`
	Duration   float64
	ViewCount  float64 `json:"view_count"`
}

var apiKey = env.MustGet("YOUTUBE_API_KEY")

func Search(ctx context.Context, query string) ([]Result, error) {
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("Error creating YouTube service: %w", err)
	}

	// Step 1: Search for videos
	searchCall := youtubeService.Search.List([]string{"id", "snippet"}).
		Q(query).       // Search query
		MaxResults(128) // Number of results to retrieve

	searchResponse, err := searchCall.Do()
	if err != nil {
		return nil, fmt.Errorf("Error making search API call: %w", err)
	}

	// Collect video IDs
	var videoIDs []string
	for _, item := range searchResponse.Items {
		if item.Id.Kind == "youtube#video" {
			videoIDs = append(videoIDs, item.Id.VideoId)
		}
	}

	// Step 2: Get video details including duration
	videoCall := youtubeService.Videos.List([]string{"id", "contentDetails", "snippet", "statistics"}).
		Id(videoIDs...) // Pass the video IDs

	videoResponse, err := videoCall.Do()
	if err != nil {
		return nil, fmt.Errorf("Error making videos API call: %w", err)
	}

	results := []Result{}

	// Parse and display results
	fmt.Println("Video Details:")
	for _, video := range videoResponse.Items {
		if video.Snippet.LiveBroadcastContent != "none" {
			continue
		}

		duration := video.ContentDetails.Duration
		parsedDuration, err := parseYouTubeDuration(duration)
		if err != nil {
			return nil, fmt.Errorf("Error parsing duration for video %s: %w", video.Id, err)
		}

		xduration, err := time.ParseDuration(parsedDuration)
		if err != nil {
			return nil, fmt.Errorf("Error parsing parsedDuration for video %s: %w", video.Id, err)
		}

		results = append(results, Result{
			ID:         video.Id,
			Thumbnail:  video.Snippet.Thumbnails.Medium.Url,
			Title:      video.Snippet.Title,
			WebpageURL: fmt.Sprintf("https://www.youtube.com/watch?v=%s", video.Id),
			Duration:   xduration.Seconds(),
			ViewCount:  float64(video.Statistics.ViewCount),
			//LikeCount:  video.Statistics.LikeCount,
			Uploader: video.Snippet.ChannelTitle,
		})
		// fmt.Printf("Title: %s\n", video.Snippet.Title)
		// fmt.Printf("Video ID: %s\n", video.Id)
		// fmt.Printf("Duration: %s\n", parsedDuration)
		// fmt.Println("----")
	}
	return results, nil
}

// parseYouTubeDuration parses the ISO 8601 duration format returned by YouTube into a more readable format
func parseYouTubeDuration(duration string) (string, error) {
	parsed, err := time.ParseDuration(duration)
	if err == nil {
		return parsed.String(), nil
	}

	// Fallback: Handle ISO 8601 duration
	// Example: "PT1H2M3S" -> "1h 2m 3s"
	var hours, minutes, seconds int
	_, err = fmt.Sscanf(duration, "PT%dH%dM%dS", &hours, &minutes, &seconds)
	if err == nil {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds), nil
	}
	_, err = fmt.Sscanf(duration, "PT%dH", &hours)
	if err == nil {
		return fmt.Sprintf("%dh", hours), nil
	}
	_, err = fmt.Sscanf(duration, "PT%dM%dS", &minutes, &seconds)
	if err == nil {
		return fmt.Sprintf("%dm%ds", minutes, seconds), nil
	}
	_, err = fmt.Sscanf(duration, "PT%dM", &minutes)
	if err == nil {
		return fmt.Sprintf("%dm", minutes), nil
	}
	_, err = fmt.Sscanf(duration, "PT%dS", &seconds)
	if err == nil {
		return fmt.Sprintf("%ds", seconds), nil
	}

	return "", fmt.Errorf("invalid duration format: %s", duration)
}
