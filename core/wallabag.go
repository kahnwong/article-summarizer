package core

import (
	"github.com/Strubbl/wallabago/v9"
)

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

func GetEntries() ([]wallabago.Item, error) {
	// get newest 5 articles
	entries, err := wallabago.GetEntries(
		wallabago.APICall,
		0, 0, "", "", 1, 5, "", 0, -1, "", "")

	return entries.Embedded.Items, err
}
