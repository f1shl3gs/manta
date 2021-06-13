package env

import (
	"os"
	"strconv"
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

func EnvOrDefaultInt(name string, def int) int {
	val := os.Getenv(name)
	if val == "" {
		return def
	}

	v, err := strconv.Atoi(name)
	if err != nil {
		panic(err)
	}

	return v
}
