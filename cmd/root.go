/*
Copyright © 2024 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	goose "github.com/advancedlogic/GoOse"

	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

type article struct {
	Title   string
	Content string
}

func extractArticle(url string) (article, error) {
	g := goose.New()
	articleData, err := g.ExtractFromURL(url)

	return article{Title: articleData.Title, Content: articleData.CleanedText}, err
}

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
		// validate input
		if len(args) == 0 {
			fmt.Println("Please specify URL")
			os.Exit(1)
		}

		// extract article
		article, err := extractArticle(args[0])
		if err != nil {
			log.Fatal(err)
		}

		// print article title
		fmt.Printf("========== %s ==========\n", article.Title)
		fmt.Println("")

		// summarize
		err = summarize(article.Content)
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
