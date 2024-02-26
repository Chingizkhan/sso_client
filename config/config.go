package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

var Cfg *Config

func init() {
	cfg, err := New()
	if err != nil {
		panic(fmt.Sprintf("can not init config: %v", err))
		return
	}
	Cfg = cfg
}

type (
	Config struct {
		App          `yaml:"app"`
		OauthService `yaml:"oauth_service"`
		Cookie       `yaml:"cookie"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	OauthService struct {
		Addr         string `env-required:"true" yaml:"addr" env:"OAUTH2_CLIENT_ID"`
		ClientID     string `env-required:"true" yaml:"client_id" env:"OAUTH_SERVICE_CLIENT_ID"`
		ClientSecret string `env-required:"true" yaml:"client_secret" env:"OAUTH_SERVICE_CLIENT_SECRET"`
	}

	Cookie struct {
		Secret string `env-required:"true" yaml:"secret" env:"COOKIE_SECRET"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
