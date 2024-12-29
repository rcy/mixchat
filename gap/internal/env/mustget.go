package env

import (
	"fmt"
	"os"
)

func MustGet(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s not set!", key))
	}
	return val
}
