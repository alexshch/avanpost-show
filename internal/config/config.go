package config

import (
	"avanpost-show/pkg/config"
)

const (
	ENV_CONFIG_PATH_KEY_SERVICE = "API_CONFIG_PATH"
	DEFAULT_CONFIG_PATH_SERVICE = "./config.yaml"
)

type Config struct {
	ServiceName            string                      `yaml:"service_name"     env:"SERVICE_NAME"`
	Db                     config.Database             `yaml:"db"`
	LogLevel               string                      `yaml:"log_level"        env:"LOG_LEVEL"`
	Http                   Http                        `yaml:"http"`
	HashiCorpVault         config.HashiCorpVaultConfig `yaml:"hashi_corp_vault"`
	Nats                   config.NatsConfig           `yaml:"nats"`
	SessionProcessorMSEURL string                      `yaml:"session_processor_mse_url" env:"SESSION_PROCESSOR_MSE_URL"`
}

type Http struct {
	Port       string `yaml:"port"        env:"API_PORT"`
	BaseAPI    string `yaml:"base_api"    env:"BASE_API"`
	SigningKey string `yaml:"signing_key" env:"SIGNING_KEY"`
}

func NewConfigLoader() config.Loader[Config] {
	return config.Loader[Config]{
		Env:         ENV_CONFIG_PATH_KEY_SERVICE,
		DefaultConf: DEFAULT_CONFIG_PATH_SERVICE,
	}
}

func (cfg *Config) Addr() string {
	return cfg.Http.Port
}
