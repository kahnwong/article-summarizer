package core

import (
	"context"
	"fmt"
	"strings"

	"charm.land/glamour/v2"
	"github.com/Strubbl/wallabago/v9"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/genai"
)

func SummarizeArticle(entry wallabago.Item, mode string) (string, error) {
	p := bluemonday.StripTagsPolicy()
	contentSanitized := p.Sanitize(entry.Content)
	return Summarize(contentSanitized, DetectLanguage(entry.Content), mode)
}

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
	prompt := fmt.Sprintf("Please summarize the text using precise and concise language. Use headers and bulleted lists in the summary, to make it scannable. Maintain the meaning and factual accuracy. %s.", content)

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
	for resp, err := range client.Models.GenerateContentStream(ctx, "gemini-3.1-flash-lite-preview", genai.Text(prompt), nil) {
		if err != nil {
			return "", fmt.Errorf("failed to generate text: %w", err)
		}

		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				output += part.Text
			}
		}
	}

	if mode == "cli" {
		rendered, err := glamour.RenderWithEnvironmentConfig(output)
		if err != nil {
			return "", fmt.Errorf("failed to render markdown: %w", err)
		}
		fmt.Print(rendered)
		return "", nil
	}

	return output, nil
}
