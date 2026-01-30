package core

import (
	cliBase "github.com/kahnwong/cli-base-sops"
	"github.com/rs/zerolog/log"
)

type Config struct {
	WallabagUrl    string `yaml:"WALLABAG_URL"`
	ClientID       string `yaml:"CLIENT_ID"`
	ClientSecret   string `yaml:"CLIENT_SECRET"`
	Username       string `yaml:"USERNAME"`
	Password       string `yaml:"PASSWORD"`
	GoogleAIApiKey string `yaml:"GOOGLE_AI_API_KEY"`
}

var AppConfig *Config

func init() {
	config, err := cliBase.ReadYamlSops[Config]("~/.config/article-summarizer/config.sops.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}
	AppConfig = config
}
