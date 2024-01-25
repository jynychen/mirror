package main

import (
	"fmt"
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
	errOccurred bool
)

func init() {
	appLogger = logger.New(os.Stdout)

	if err := yaml.Unmarshal([]byte(os.Getenv(config.CONFIG_ENV)), &appConfig); err != nil {
		appLogger.Fatal("Error config yaml.Unmarshal():", "err", err)
	}
	if len(appConfig.MirrorConfigs) == 0 {
		appLogger.Fatal("no mirror config found")
	}

	for _, cfg := range appConfig.MirrorConfigs {
		mirrorTasks = append(mirrorTasks, mirror.New(&mirror.MirrorConfig{
			Logger:     appLogger,
			SrcRepoURL: os.Getenv(cfg.EnvSrcRepoURL),
			SrcAuth:    config.NewNewPublicKeysOrNil(os.Getenv(cfg.EnvSrcSShKey)),
			DstRepoURL: os.Getenv(cfg.EnvDstRepoURL),
			DstAuth:    config.NewNewPublicKeysOrNil(os.Getenv(cfg.EnvDstSShKey)),
		}))
	}
}

func main() {
	for _, mirrorTask := range mirrorTasks {
		fmt.Println()
		if err := mirrorTask.Run(); err != nil {
			appLogger.Warn("Error mirrorTask.Run():", "err", err)
			errOccurred = true
		}
		fmt.Println()
	}

	if errOccurred {
		os.Exit(1)
	}
}
