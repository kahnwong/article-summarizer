package cmd

import (
	"github.com/Strubbl/wallabago/v9"
)

type WallabagConfig struct {
	WallabagUrl  string `yaml:"WALLABAG_URL"`
	ClientID     string `yaml:"CLIENT_ID"`
	ClientSecret string `yaml:"CLIENT_SECRET"`
	Username     string `yaml:"USERNAME"`
	Password     string `yaml:"PASSWORD"`
}

func init() {
	config := readConfig()

	wallabagConfig := wallabago.WallabagConfig{
		WallabagURL:  config.WallabagUrl,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		UserName:     config.Username,
		UserPassword: config.Password,
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
