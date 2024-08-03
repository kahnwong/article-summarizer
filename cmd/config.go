package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/getsops/sops/v3/decrypt"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func readConfig() WallabagConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err)
	}
	filename := filepath.Join(homeDir, ".config", "article-summarizer", "config.sops.yaml")

	// Check if the file exists
	_, err = os.Stat(filename)

	if os.IsNotExist(err) {
		log.Panic().Msgf("File %s does not exist.", filename)
		os.Exit(1)
	}

	var config WallabagConfig

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
