/*
Copyright © 2024 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

func summarize(content string) error {
	// prep
	ollamaModel := "kahnwong/gemma-1.1:7b-it"

	//if language == "th":
	//prompt += "Respond in Thai language."
	//model_name = "kahnwong/typhoon-1.5:8b"

	// init ollama
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	// ollama request payload
	req := &api.GenerateRequest{
		Model:  ollamaModel,
		Prompt: fmt.Sprintf("summarize following text into four paragraphs: %s.", content),
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
		log.Fatal(err)
	}

	return err
}

var rootCmd = &cobra.Command{
	Use:   "article-summarizer",
	Short: "Summarize an article with LLM",
	Run: func(cmd *cobra.Command, args []string) {
		err := summarize("I'm a swimmer from a faraway land")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
