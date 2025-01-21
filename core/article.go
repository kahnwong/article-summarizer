package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
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
	client, err := genai.NewClient(ctx, option.WithAPIKey(AppConfig.GoogleAIApiKey))
	if err != nil {
		log.Fatal().Msg("Failed to create GOOGLE AI client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	iter := model.GenerateContentStream(ctx, genai.Text(prompt))

	var output string
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal().Msg("Failed to generate text")
		}

		if resp.Candidates != nil {
			for _, v := range resp.Candidates {
				for _, k := range v.Content.Parts {
					if mode == "cli" {
						time.Sleep(1 * time.Second)
						fmt.Print(k.(genai.Text))
					} else {
						output += fmt.Sprint(k.(genai.Text))
					}
				}
			}
		}
	}

	return output
}
