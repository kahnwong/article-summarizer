package core

import (
	"context"
	"fmt"
	"strings"

	"charm.land/glamour/v2"
	"github.com/Strubbl/wallabago/v9"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/interactions"
	"google.golang.org/genai/interactions/models/operations"
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
		APIKey:  AppConfig.GoogleApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create GOOGLE AI client: %w", err)
	}

	res, err := client.Interactions.Create(ctx, operations.CreateInteractionRequest{
		Body: operations.NewCreateInteractionRequestBody(interactions.CreateModelInteraction{
			Model:  interactions.Model(AppConfig.ModelName),
			Input:  new(interactions.NewInteractionsInput(prompt)),
			Stream: new(true),
		}),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create interaction: %w", err)
	}

	eventStream := res.InteractionSSEStreamEvent
	if eventStream == nil {
		return "", fmt.Errorf("failed to create interaction stream")
	}
	defer eventStream.Close()

	var output strings.Builder
	for eventStream.Next() {
		textDelta := eventStream.Value().GetDataStepDelta().GetDeltaText()
		if textDelta != nil {
			output.WriteString(textDelta.Text)
		}
	}
	if err := eventStream.Err(); err != nil {
		return "", fmt.Errorf("failed to generate text: %w", err)
	}

	if mode == "cli" {
		rendered, err := glamour.RenderWithEnvironmentConfig(output.String())
		if err != nil {
			return "", fmt.Errorf("failed to render markdown: %w", err)
		}
		fmt.Print(rendered)
		return "", nil
	}

	return output.String(), nil
}
