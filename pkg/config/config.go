package config

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
