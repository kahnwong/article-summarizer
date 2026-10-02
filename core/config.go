package core

import (
	"fmt"

	cliBase "github.com/kahnwong/cli-base-sops"
)

type Config struct {
	WallabagUrl  string `yaml:"WALLABAG_URL"`
	ClientID     string `yaml:"CLIENT_ID"`
	ClientSecret string `yaml:"CLIENT_SECRET"`
	Username     string `yaml:"USERNAME"`
	Password     string `yaml:"PASSWORD"`
	GoogleApiKey string `yaml:"GOOGLE_API_KEY"`
	ModelName    string `yaml:"MODEL_NAME"`
}

var AppConfig *Config

func LoadConfig() error {
	config, err := cliBase.ReadYamlSops[Config]("~/.config/article-summarizer/config.sops.yaml")
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	AppConfig = config
	return nil
}
