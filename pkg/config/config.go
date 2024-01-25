package config

import (
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	xssh "golang.org/x/crypto/ssh"
)

const (
	CONFIG_ENV = "CONFIG"
)

type MirrorConfigs []MirrorConfig

type MirrorConfig struct {
	EnvSrcRepoURL string `yaml:"env_src_repo_url"`
	EnvSrcSShKey  string `yaml:"env_src_ssh_key"`
	EnvDstRepoURL string `yaml:"env_dst_repo_url"`
	EnvDstSShKey  string `yaml:"env_dst_ssh_key"`
}

type Config struct {
	MirrorConfigs MirrorConfigs `yaml:"mirror_configs"`
}

func NewNewPublicKeysOrNil(pemBytes string) ssh.AuthMethod {
	pk, err := ssh.NewPublicKeys("git", []byte(pemBytes), "")
	if err != nil {
		return nil
	}

	pk.HostKeyCallbackHelper.HostKeyCallback = xssh.InsecureIgnoreHostKey()
	return pk
}
