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
	srcRepoURL string
	srcAuth    transport.AuthMethod
	dstRepoURL string
	dstAuth    transport.AuthMethod
	logger     logger.Logger
}

type MirrorConfig struct {
	SrcRepoURL string
	SrcAuth    transport.AuthMethod
	DstRepoURL string
	DstAuth    transport.AuthMethod
	Logger     logger.Logger
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
	m.logger.Info("Start mirroring...")
	defer m.logger.Info("End mirroring...")

	if m.srcRepoURL == "" {
		m.logger.Error("validate src_repo_url", "err", ErrEmptySrcRepoURL)
		return ErrEmptySrcRepoURL
	}
	m.logger.Debug("Source:", "repo", m.srcRepoURL)

	if m.dstRepoURL == "" {
		m.logger.Error("validate dst_repo_url", "err", ErrEmptyDstRepoURL)
		return ErrEmptyDstRepoURL
	}
	m.logger.Debug("Destination:", "repo", m.dstRepoURL)

	m.logger.Info("Cloning source repository...")
	srcRepo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:      m.srcRepoURL,
		Auth:     m.srcAuth,
		Progress: m.logger,
		Mirror:   true,
	})
	if err != nil {
		m.logger.Error("failed to clone source repository", "err", err)
		return err
	}

	_, err = srcRepo.CreateRemote(&config.RemoteConfig{
		Name: MirrorRemote,
		URLs: []string{m.dstRepoURL},
	})
	if err != nil {
		m.logger.Error("failed to create remote", "err", err)
		return err
	}

	m.logger.Info("Pushing to destination repository...")
	err = srcRepo.Push(&git.PushOptions{
		RemoteName: MirrorRemote,
		Auth:       m.dstAuth,
		Progress:   m.logger,
		RefSpecs: []config.RefSpec{
			"+refs/*:refs/*",
		},
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		m.logger.Error("failed to push", "err", err)
		return err
	}

	m.logger.Debug("Successfully mirrored.")
	return nil
}
