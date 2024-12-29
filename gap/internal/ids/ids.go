package ids

import "github.com/oklog/ulid/v2"

func Make(prefix string) string {
	return prefix + "_" + ulid.Make().String()
}

func MakeTrackID() string {
	return Make("trk")
}
