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

	"github.com/microcosm-cc/bluemonday"

	"github.com/ollama/ollama/api"

	"github.com/Strubbl/wallabago/v9"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func getEntries() ([]wallabago.Item, error) {
	// get newest 5 articles
	entries, err := wallabago.GetEntries(
		wallabago.APICall,
		0, 0, "", "", 1, 5, "", 0, -1, "", "")
	if err != nil {
		return nil, err
	}

	return entries.Embedded.Items, err
}

func createFormOptions(entries []wallabago.Item) []huh.Option[string] {
	var options []huh.Option[string]

	for _, v := range entries {
		options = append(options, huh.NewOption(v.Title, v.Title).Selected(true))

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

		// ------------ get entries ------------ //
		wallabago.SetConfig(wallabagConfig)

		entries, err := getEntries()
		if err != nil {
			log.Println("Cannot obtain articles from Wallabag")
		}

		var (
			entryTitle string
		)

		// ------------ select article ------------ //
		formEntries := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose an article to summarize").
					Options(
						createFormOptions(entries)...,
					).
					Value(&entryTitle),
			),
		)
		err = formEntries.Run()
		if err != nil {
			log.Fatal(err)
		}

		// ------------ summarize ------------ //
		fmt.Printf("========== %s ==========\n", entryTitle)

		var content string
		for _, entry := range entries {
			if entry.Title == entryTitle {
				content = entry.Content
			}
		}

		p := bluemonday.StripTagsPolicy()
		contentSanitized := p.Sanitize(
			content,
		)

		err = summarize(contentSanitized, detectLanguage(content))
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
