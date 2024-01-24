package config

import "github.com/go-git/go-git/v5/plumbing/transport/ssh"

const (
	CONFIG_ENV = "CONFIG"
)

type MirrorConfigs []MirrorConfig

type MirrorConfig struct {
	SrcRepoURL string `yaml:"src_repo_url"`
	SrcSShKey  string `yaml:"src_ssh_key"`
	DstRepoURL string `yaml:"dst_repo_url"`
	DstSShKey  string `yaml:"dst_ssh_key"`
}

type Config struct {
	MirrorConfigs MirrorConfigs `yaml:"mirror_configs"`
}

func NewNewPublicKeysOrNil(pemBytes string) ssh.AuthMethod {
	pk, err := ssh.NewPublicKeys("git", []byte(pemBytes), "")
	if err != nil {
		return nil
	}

	return pk
}
