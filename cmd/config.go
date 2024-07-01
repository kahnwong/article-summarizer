package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Strubbl/wallabago/v9"

	"github.com/getsops/sops/v3/decrypt"
	"gopkg.in/yaml.v3"
)

var wallabagConfig = createWallabagConfig()

type Config struct {
	WallabagUrl  string `yaml:"WALLABAG_URL"`
	ClientID     string `yaml:"CLIENT_ID"`
	ClientSecret string `yaml:"CLIENT_SECRET"`
	Username     string `yaml:"USERNAME"`
	Password     string `yaml:"PASSWORD"`
}

func readConfig() Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	filename := filepath.Join(homeDir, ".config", "article-summarizer", "config.sops.yaml")

	// Check if the file exists
	_, err = os.Stat(filename)

	if os.IsNotExist(err) {
		fmt.Printf("File %s does not exist.\n", filename)
		os.Exit(1)
	}

	var config Config

	data, err := decrypt.File(filename, "yaml")
	if err != nil {
		fmt.Println(fmt.Printf("Failed to decrypt: %v", err))
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return config
}

func createWallabagConfig() wallabago.WallabagConfig {
	config := readConfig()

	return wallabago.WallabagConfig{
		WallabagURL:  config.WallabagUrl,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		UserName:     config.Username,
		UserPassword: config.Password,
	}
}
