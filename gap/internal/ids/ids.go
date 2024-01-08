package ids

import "github.com/oklog/ulid/v2"

func Make(prefix string) string {
	return prefix + "_" + ulid.Make().String()
}
