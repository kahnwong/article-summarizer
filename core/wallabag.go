package core

import (
	"fmt"
	"strconv"

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
	entries, err := wallabago.GetEntries(
		wallabago.APICall,
		0, 0, "", "", 1, 5, "", 0, -1, "", "")

	return entries.Embedded.Items, err
}

func MarkEntryAsRead(id int) error {
	url := AppConfig.WallabagUrl + "/api/entries/" + strconv.Itoa(id) + ".json"
	_, err := wallabago.APICall(url, "PATCH", []byte(`{"archive": 1}`))
	if err != nil {
		return fmt.Errorf("failed to mark entry as read: %w", err)
	}
	return nil
}
