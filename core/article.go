package core

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func Summarize(content string, language string, mode string) (string, error) {
	// set parameters
	prompt := fmt.Sprintf("summarize following text into four paragraphs: %s.", content)

	if language == "Thai" {
		prompt += "Respond in Thai language."
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  AppConfig.GoogleAIApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create GOOGLE AI client: %w", err)
	}

	var output string
	for resp, err := range client.Models.GenerateContentStream(ctx, "gemini-3-flash-preview", genai.Text(prompt), nil) {
		if err != nil {
			return "", fmt.Errorf("failed to generate text: %w", err)
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

	return output, nil
}
