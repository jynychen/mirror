package mirror

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jyny/mirror/pkg/logger"
)

const (
	MirrorRemote = "mirror"
)

type Mirror struct {
	srcRepoURL string
	srcAuth    transport.AuthMethod
	dstRepoURL string
	dstAuth    transport.AuthMethod
	logger     *logger.Logger
}

type MirrorConfig struct {
	SrcRepoURL string
	SrcAuth    transport.AuthMethod
	DstRepoURL string
	DstAuth    transport.AuthMethod
	Logger     *logger.Logger
}

func New(cfg *MirrorConfig) *Mirror {
	return &Mirror{
		srcRepoURL: cfg.SrcRepoURL,
		srcAuth:    cfg.SrcAuth,
		dstRepoURL: cfg.DstRepoURL,
		dstAuth:    cfg.DstAuth,
		logger:     cfg.Logger,
	}
}

func (m *Mirror) Run() error {
	m.logger.Println("Start mirroring...")
	defer m.logger.Printf("End mirroring...\n\n")

	if m.srcRepoURL == "" {
		return ErrEmptySrcRepoURL
	}
	m.logger.Println("Source Repository: ", m.srcRepoURL)

	if m.dstRepoURL == "" {
		return ErrEmptyDstRepoURL
	}
	m.logger.Println("Destination Repository: ", m.dstRepoURL)

	m.logger.Println("Cloning source repository...")
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
		URLs: []string{m.dstRepoURL},
	})
	if err != nil {
		return err
	}

	m.logger.Println("Pushing to destination repository...")
	err = srcRepo.Push(&git.PushOptions{
		RemoteName: MirrorRemote,
		Auth:       m.dstAuth,
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
