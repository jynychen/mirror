package main

import (
	"github.com/jynychen/mirror/pkg/config"
	"github.com/jynychen/mirror/pkg/logger"
	"github.com/jynychen/mirror/pkg/mirror"
)

var (
	appLogger *logger.Logger
	appConfig config.Config
)

func init() {
	appLogger = logger.New()
}

func main() {
	cfg := &mirror.MirrorConfig{
		Logger: appLogger,
	}

	if err := mirror.New(cfg).Run(); err != nil {
		appLogger.Fatalf("Mirror.Run() err: %v", err)
	}
}
