package main

import (
	"github.com/jynychen/mirror/pkg/logger"
	"github.com/jynychen/mirror/pkg/mirror"
)

var appLogger *logger.Logger

func init() {
	appLogger = logger.New()
}

func main() {
	srcRepoURL := ""
	destRepoURL := ""

	if err := mirror.New(srcRepoURL, destRepoURL, appLogger).Run(); err != nil {
		appLogger.Fatalf("Mirror.Run() err: %v", err)
	}
}
