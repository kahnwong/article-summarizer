package cmd

import (
	"github.com/Strubbl/wallabago/v9"
)

type Config struct {
	WallabagUrl    string `yaml:"WALLABAG_URL"`
	ClientID       string `yaml:"CLIENT_ID"`
	ClientSecret   string `yaml:"CLIENT_SECRET"`
	Username       string `yaml:"USERNAME"`
	Password       string `yaml:"PASSWORD"`
	GoogleAIApiKey string `yaml:"GOOGLE_AI_API_KEY"`
}

func init() {
	wallabagConfig := wallabago.WallabagConfig{
		WallabagURL:  AppConfig.WallabagUrl,
		ClientID:     AppConfig.ClientID,
		ClientSecret: AppConfig.ClientSecret,
		UserName:     AppConfig.Username,
		UserPassword: AppConfig.Password,
	}
	wallabago.SetConfig(wallabagConfig)
}

func getEntries() ([]wallabago.Item, error) {
	// get newest 5 articles
	entries, err := wallabago.GetEntries(
		wallabago.APICall,
		0, 0, "", "", 1, 5, "", 0, -1, "", "")

	return entries.Embedded.Items, err
}
