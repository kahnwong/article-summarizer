package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/Strubbl/wallabago/v9"
	"github.com/charmbracelet/huh"
	"github.com/ollama/ollama/api"
	"github.com/rs/zerolog/log"
)

func createFormOptions(entries []wallabago.Item) []huh.Option[string] {
	var options []huh.Option[string]

	for _, v := range entries {
		options = append(options, huh.NewOption(v.Title, v.Title))
	}

	return options
}

func detectLanguage(content string) string {
	var language string
	if !strings.Contains(content, "ก") {
		language = "English"
	} else {
		language = "Thai"
	}

	return language
}

func summarize(content string, language string) {
	// set parameters
	ollamaModel := "kahnwong/gemma-1.1:7b-it"
	prompt := fmt.Sprintf("summarize following text into four paragraphs: %s.", content)

	if language == "Thai" {
		prompt += "Respond in Thai language."
		ollamaModel = "kahnwong/typhoon-1.5:8b"
	}

	// init ollama
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Error().Msg("Could not init ollama client")
	}

	// ollama request payload
	req := &api.GenerateRequest{
		Model:  ollamaModel,
		Prompt: prompt,
	}

	// render results
	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response)
		return nil
	}

	// main
	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Error().Msg("Could not summarize article")
	}
}
