/*
Copyright © 2024 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Strubbl/wallabago/v9"
	"github.com/spf13/cobra"

	"github.com/rs/zerolog/log"

	"github.com/microcosm-cc/bluemonday"

	"github.com/charmbracelet/huh"
	"github.com/kahnwong/article-summarizer/core"
)

var entryTitle string // for huh form

// functions
func createFormOptions(entries []wallabago.Item) []huh.Option[string] {
	var options []huh.Option[string]

	for _, v := range entries {
		options = append(options, huh.NewOption(v.Title, v.Title))
	}

	return options
}

// main
var rootCmd = &cobra.Command{
	Use:   "article-summarizer",
	Short: "Summarize an article with LLM",
	Run: func(cmd *cobra.Command, args []string) {
		// Clears the screen
		core.ClearScreen()

		// ------------ get entries ------------ //
		entries, err := core.GetEntries()
		if err != nil {
			log.Fatal().Msg("Cannot obtain articles from Wallabag")
		}

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
			log.Fatal().Err(err)
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

		core.Summarize(contentSanitized, core.DetectLanguage(content), "cli")
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
