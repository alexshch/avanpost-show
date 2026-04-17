package config

type HashiCorpVaultConfig struct {
	ServerAddr         string `yaml:"server_addr" env:"HASHI_CORP_SERVER_ADDR" json:"server_addr,omitempty"`
	RoleID             string `yaml:"role_id" env:"HASHI_CORP_ROLE_ID" json:"role_id,omitempty"`
	SecretID           string `yaml:"secret_id" env:"HASHI_CORP_SECRET_ID" json:"secret_id,omitempty"`
	MountPath          string `yaml:"mount_path" env:"HASHI_CORP_MOUNT_PATH" json:"mount_path,omitempty"`
	KeyPrefix          string `yaml:"key_prefix" env:"HASHI_CORP_KEY_PREFIX" json:"key_prefix,omitempty"`
	Tls                bool   `yaml:"tls" env:"HASHI_CORP_USE_TLS" json:"tls,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify" env:"HASHI_CORP_INSECURE_SKIP_VERIFY" json:"insecure_skip_verify,omitempty" env-default:"true"`
}
