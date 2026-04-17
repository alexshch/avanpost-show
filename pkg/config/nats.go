package config

type NatsConfig struct {
	Url      string `yaml:"url"      env:"NATS_URL"      env-default:"nats://127.0.0.1:4222"`
	Login    string `yaml:"login"    env:"NATS_LOGIN"    env-default:"nats_login"`
	Password string `yaml:"password" env:"NATS_PASSWORD" env-default:"nats_password"`
	Replicas int    `yaml:"replicas" env:"NATS_REPLICAS" env-default:"1"`
}
