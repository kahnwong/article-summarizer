/*
Copyright © 2024 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	goose "github.com/advancedlogic/GoOse"
	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

type article struct {
	Title    string
	Content  string
	Language string
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

func extractArticle(url string) (article, error) {
	g := goose.New()
	articleData, err := g.ExtractFromURL(url)

	// detect language
	language := detectLanguage(articleData.CleanedText)

	return article{Title: articleData.Title, Content: articleData.CleanedText, Language: language}, err
}

func summarize(content string, language string) error {
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
		log.Fatal(err)
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
		log.Fatal(err)
	}

	return err
}

var rootCmd = &cobra.Command{
	Use:   "article-summarizer",
	Short: "Summarize an article with LLM",
	Run: func(cmd *cobra.Command, args []string) {
		// Clears the screen
		ClearScreen()

		// validate input
		var url string
		if len(args) == 0 {
			urlFromClipboard := ReadFromClipboard()
			if urlFromClipboard != "" {
				if strings.HasPrefix(urlFromClipboard, "https://") {
					url = urlFromClipboard
				}
			}
		}
		if url == "" {
			if len(args) == 0 {
				fmt.Println("Please specify URL")
				os.Exit(1)
			} else if len(args) == 1 {
				url = args[0]
			}
		}

		// extract article
		article, err := extractArticle(url)
		if err != nil {
			log.Fatal(err)
		}
		if article.Content == "" {
			fmt.Println("Could not extract article")
			os.Exit(1)
		}

		// print article title
		fmt.Printf("========== %s ==========\n", article.Title)
		fmt.Println("")

		// summarize
		err = summarize(article.Content, article.Language)
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
