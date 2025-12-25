package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/genai"
)

func DetectLanguage(content string) string {
	var language string
	if !strings.Contains(content, "ก") {
		language = "English"
	} else {
		language = "Thai"
	}

	return language
}

func Summarize(content string, language string, mode string) string {
	// set parameters
	//ollamaModel := "kahnwong/gemma-1.1:7b-it"
	prompt := fmt.Sprintf("summarize following text into four paragraphs: %s.", content)

	if language == "Thai" {
		prompt += "Respond in Thai language."
		//	ollamaModel = "kahnwong/typhoon-1.5:8b"
	}

	//// init ollama
	//client, err := api.ClientFromEnvironment()
	//if err != nil {
	//	log.Fatal().Msg("Could not init ollama client")
	//}
	//
	//// ollama request payload
	//req := &api.GenerateRequest{
	//	Model:  ollamaModel,
	//	Prompt: prompt,
	//}
	//
	//// render results
	//ctx := context.Background()
	//respFunc := func(resp api.GenerateResponse) error {
	//	fmt.Print(resp.Response)
	//	return nil
	//}
	//
	//// main
	//err = client.Generate(ctx, req, respFunc)
	//if err != nil {
	//	log.Error().Msg("Could not summarize article")
	//}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  AppConfig.GoogleAIApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal().Msg("Failed to create GOOGLE AI client")
	}

	var output string
	for resp, err := range client.Models.GenerateContentStream(ctx, "gemini-3-flash-preview", genai.Text(prompt), nil) {
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to generate text")
		}

		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				text := part.Text
				if mode == "cli" {
					time.Sleep(1 * time.Second)
					fmt.Print(text)
				} else {
					output += text
				}
			}
		}
	}

	return output
}
