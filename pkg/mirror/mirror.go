package mirror

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jynychen/mirror/pkg/logger"
)

const (
	MirrorRemote = "mirror"
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
	defer m.logger.Printf("End mirroring...\n\n")

	if m.srcRepoURL == "" {
		return ErrEmptySourceRepoURL
	}
	m.logger.Println("Source repo: ", m.srcRepoURL)

	if m.destRepoURL == "" {
		return ErrEmptyDestinationURL
	}
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

	_, err = srcRepo.CreateRemote(&config.RemoteConfig{
		Name: MirrorRemote,
		URLs: []string{m.destRepoURL},
	})
	if err != nil {
		return err
	}

	m.logger.Println("Pushing to destination repo...")
	err = srcRepo.Push(&git.PushOptions{
		RemoteName: MirrorRemote,
		Auth:       m.destAuth,
		Progress:   m.logger,
		RefSpecs: []config.RefSpec{
			"+refs/*:refs/*",
		},
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	m.logger.Println("Successfully mirrored.")
	return nil
}
