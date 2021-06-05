package env

import (
	"os"
	"time"
)

func EnvOrDefaultDuration(name string, def time.Duration) time.Duration {
	val := os.Getenv(name)
	if val == "" {
		return def
	}

	v, err := time.ParseDuration(val)
	if err != nil {
		panic(err)
	}

	return v
}
