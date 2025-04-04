// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Event struct {
	EventID   string
	EventType string
	CreatedAt pgtype.Timestamptz
	Payload   []byte
}

type Result struct {
	ResultID  string
	SearchID  string
	StationID string
	CreatedAt pgtype.Timestamptz
	ExternID  string
	URL       string
	Thumbnail string
	Title     string
	Uploader  string
	Duration  float64
	Views     float64
}

type RiverClient struct {
	ID        string
	CreatedAt pgtype.Timestamptz
	Metadata  []byte
	PausedAt  pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type RiverClientQueue struct {
	RiverClientID    string
	Name             string
	CreatedAt        pgtype.Timestamptz
	MaxWorkers       int64
	Metadata         []byte
	NumJobsCompleted int64
	NumJobsRunning   int64
	UpdatedAt        pgtype.Timestamptz
}

type RiverJob struct {
	ID           int64
	State        interface{}
	Attempt      int16
	MaxAttempts  int16
	AttemptedAt  pgtype.Timestamptz
	CreatedAt    pgtype.Timestamptz
	FinalizedAt  pgtype.Timestamptz
	ScheduledAt  pgtype.Timestamptz
	Priority     int16
	Args         []byte
	AttemptedBy  []string
	Errors       [][]byte
	Kind         string
	Metadata     []byte
	Queue        string
	Tags         []string
	UniqueKey    []byte
	UniqueStates pgtype.Bits
}

type RiverLeader struct {
	ElectedAt pgtype.Timestamptz
	ExpiresAt pgtype.Timestamptz
	LeaderID  string
	Name      string
}

type RiverQueue struct {
	Name      string
	CreatedAt pgtype.Timestamptz
	Metadata  []byte
	PausedAt  pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type SchemaVersion struct {
	Version int32
}

type Search struct {
	SearchID  string
	StationID string
	CreatedAt pgtype.Timestamptz
	Query     string
	Status    string
}

type Session struct {
	SessionID string
	CreatedAt pgtype.Timestamptz
	ExpiresAt pgtype.Timestamptz
	UserID    string
}

type Station struct {
	StationID          string
	CreatedAt          pgtype.Timestamptz
	Slug               string
	Name               string
	Active             bool
	CurrentTrackID     pgtype.Text
	BackgroundImageURL string
	UserID             string
	IsPublic           bool
	TelnetPort         string
	BroadcastPort      string
}

type StationMessage struct {
	StationMessageID string
	CreatedAt        pgtype.Timestamptz
	Type             string
	StationID        string
	ParentID         string
	Nick             string
	Body             string
	IsHidden         pgtype.Bool
}

type Track struct {
	TrackID     string
	StationID   string
	CreatedAt   pgtype.Timestamptz
	Artist      string
	Title       string
	RawMetadata []byte
	Rotation    int32
	Plays       int32
	Skips       int32
	Playing     bool
}

type User struct {
	UserID    string
	CreatedAt pgtype.Timestamptz
	Username  string
	Guest     bool
}
