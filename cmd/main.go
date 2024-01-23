package main

import (
	"log"
	"os"

	"github.com/jyny/mirror/pkg/config"
	"github.com/jyny/mirror/pkg/logger"
	"github.com/jyny/mirror/pkg/mirror"
	"gopkg.in/yaml.v3"
)

var (
	appConfig   config.Config
	appLogger   logger.Logger
	mirrorTasks []*mirror.Mirror
)

func init() {
	if err := yaml.Unmarshal([]byte(os.Getenv(config.CONFIG_ENV)), &appConfig); err != nil {
		log.Fatalf("Error config yaml.Unmarshal(): %v", err)
	}
	if len(appConfig.MirrorConfigs) == 0 {
		log.Fatalf("no mirror config found")
	}

	appLogger = logger.New()
}

func main() {
	for _, each := range appConfig.MirrorConfigs {
		mirrorTasks = append(mirrorTasks, mirror.New(&mirror.MirrorConfig{
			Logger:     appLogger,
			SrcRepoURL: each.SrcRepoURL,
			DstRepoURL: each.DstRepoURL,
		}))
	}

	for _, each := range mirrorTasks {
		if err := each.Run(); err != nil {
			log.Printf("Error mirrorTask Run(): %v\n\n", err)
		}
	}
}
