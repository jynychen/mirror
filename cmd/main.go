package main

import (
	"github.com/jynychen/mirror/pkg/log"
	"github.com/jynychen/mirror/pkg/mirror"
)

func main() {
	logger := log.New()

	srcRepoURL := ""
	destRepoURL := ""

	if err := mirror.New(srcRepoURL, destRepoURL, logger).Run(); err != nil {
		logger.Fatalf("Mirror.Run() err: %v", err)
	}
}
