package main

import (
	"fmt"
	"os"

	"github.com/jynychen/mirror/pkg/config"
	"github.com/jynychen/mirror/pkg/logger"
	"github.com/jynychen/mirror/pkg/mirror"
	"gopkg.in/yaml.v3"
)

var (
	appConfig   config.Config
	appLogger   logger.Logger
	mirrorTasks []*mirror.Mirror
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
			SrcRepoURL: cfg.SrcRepoURL,
			SrcAuth:    config.NewNewPublicKeysOrNil(cfg.SrcSShKey),
			DstRepoURL: cfg.DstRepoURL,
			DstAuth:    config.NewNewPublicKeysOrNil(cfg.DstSShKey),
		}))
	}
}

func main() {
	for _, mirrorTask := range mirrorTasks {
		fmt.Println()
		if err := mirrorTask.Run(); err != nil {
			appLogger.Warn("Error mirrorTask.Run():", "err", err)
		}
		fmt.Println()
	}
}
