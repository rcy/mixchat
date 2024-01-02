package ytdlp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dhowden/tag"
)

type format string

const (
	OGG format = "ogg"
)

type AudioTrack struct {
	Format   format
	Data     []byte
	Metadata tag.Metadata
}

func AudioTrackFromURL(ctx context.Context, url string) (*AudioTrack, error) {
	executable := "./yt/yt-dlp"

	dir, err := os.MkdirTemp("", "mixchat-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	fmt.Printf("ripping %s into %s\n", url, dir)

	destination := fmt.Sprintf("%s/track", dir)

	args := []string{
		"--keep-video",
		"--extract-audio",
		//"--audio-quality=0",
		"--audio-format=vorbis",
		"--embed-metadata",
		"--max-downloads=10",
		"--output=" + destination,
		url,
	}

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
	WebpageURL string `json:"webpage_url"`
	Duration   int32
	ViewCount  int64 `json:"view_count"`
}

func Search(ctx context.Context, query string) ([]Result, error) {
	executable := "./yt/yt-dlp"

	dir, err := os.MkdirTemp("", "mixchat-search-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	args := []string{
		"--no-download",
		"--write-info-json",
		"--output=" + "infojson:" + dir + "/tracks/%(playlist_index)s",
		"--output=" + "pl_infojson:" + dir + "/playlist",
		fmt.Sprintf("ytsearch%d:%s", 5, query),
	}

	output, err := exec.CommandContext(ctx, executable, args...).CombinedOutput()
	fmt.Println("searchoutput:", string(output))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, string(output))
	}

	entries, err := os.ReadDir(filepath.Join(dir, "tracks"))
	if err != nil {
		return nil, err
	}

	results := []Result{}
	for _, ent := range entries {
		if ent.Type().IsRegular() {
			b, err := os.ReadFile(filepath.Join(dir, "tracks", ent.Name()))
			if err != nil {
				return nil, err
			}
			var result Result
			err = json.Unmarshal(b, &result)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	}

	return results, nil
}
