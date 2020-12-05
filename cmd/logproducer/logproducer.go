package main

import (
	"io"
	"math/rand"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/f1shl3gs/manta/log"
)

func main() {
	lcf := log.Config{
		Level: zapcore.InfoLevel,
	}

	output := os.Getenv("OUTPUT")
	if output == "" {
		output = "demo.log"
	}

	f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var w io.Writer = f
	if os.Getenv("TEE_STDOUT") == "TRUE" {
		w = io.MultiWriter(f, os.Stdout)
	}

	logger, err := lcf.New(w)
	if err != nil {
		panic(err)
	}

	for {
		d := time.Duration(rand.Intn(2000)) * time.Millisecond
		time.Sleep(d)

		logger.Info("serve http request",
			zap.String("method", "GET"),
			zap.String("url", "/api/v1/demo"),
			zap.Int("status", 200),
			zap.Duration("latency", d))
	}
}
