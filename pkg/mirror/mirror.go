package mirror

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jynychen/mirror/pkg/logger"
)

type Mirror struct {
	srcRepoURL  string
	srcAuth     transport.AuthMethod
	destRepoURL string
	destAuth    transport.AuthMethod
	logger      *logger.Logger
}

type MirrorConfig struct {
	SrcRepoURL  string
	SrcAuth     transport.AuthMethod
	DestRepoURL string
	DestAuth    transport.AuthMethod
	Logger      *logger.Logger
}

func New(cfg *MirrorConfig) *Mirror {
	return &Mirror{
		srcRepoURL:  cfg.SrcRepoURL,
		srcAuth:     cfg.SrcAuth,
		destRepoURL: cfg.DestRepoURL,
		destAuth:    cfg.DestAuth,
		logger:      cfg.Logger,
	}
}

func (m *Mirror) Run() error {
	m.logger.Println("Start mirroring...")
	m.logger.Println("Source repo: ", m.srcRepoURL)
	m.logger.Println("Destination repo: ", m.destRepoURL)

	m.logger.Println("Cloning source repo...")
	srcRepo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:      m.srcRepoURL,
		Auth:     m.srcAuth,
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
		Auth:       m.destAuth,
		RefSpecs:   []config.RefSpec{"+refs/*:refs/*"},
		Progress:   m.logger,
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	m.logger.Println("Done!\n\n")
	return nil
}
