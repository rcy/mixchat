package chat

import "time"

type Message struct {
	CreatedAt   time.Time
	Nick        string
	Body        string
	StationSlug string
}

var store = make(map[string][]Message)

func Insert(stationSlug string, msg Message) {
	store[stationSlug] = append(store[stationSlug], msg)
}

func Fetch(stationSlug string) []Message {
	return store[stationSlug]
}
