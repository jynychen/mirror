package mirror

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jynychen/mirror/pkg/logger"
)

type Mirror struct {
	srcRepoURL  string
	destRepoURL string
	logger      *logger.Logger
}

func New(srcRepoURL, destRepoURL string, logger *logger.Logger) *Mirror {
	return &Mirror{
		srcRepoURL:  srcRepoURL,
		destRepoURL: destRepoURL,
		logger:      logger,
	}
}

func (m *Mirror) Run() error {
	m.logger.Println("Start mirroring...")
	m.logger.Println("Source repo: ", m.srcRepoURL)
	m.logger.Println("Destination repo: ", m.destRepoURL)

	m.logger.Println("Cloning source repo...")
	srcRepo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:      m.srcRepoURL,
		Progress: m.logger,
		Mirror:   true,
	})
	if err != nil {
		return err
	}

	mirror, err := srcRepo.CreateRemote(&config.RemoteConfig{
		Name: "mirror",
		URLs: []string{m.destRepoURL},
	})
	if err != nil {
		return err
	}

	m.logger.Println("Pushing to destination repo...")
	err = srcRepo.Push(&git.PushOptions{
		RemoteName: mirror.Config().Name,
		RefSpecs:   []config.RefSpec{"+refs/*:refs/*"},
		Progress:   m.logger,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	m.logger.Println("Done!\n\n")
	return nil
}
