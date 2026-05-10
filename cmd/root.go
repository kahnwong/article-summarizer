/*
Copyright © 2026 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Strubbl/wallabago/v9"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"github.com/spf13/cobra"

	"github.com/microcosm-cc/bluemonday"

	"charm.land/huh/v2"
	"github.com/kahnwong/article-summarizer/core"
)

var entryTitle string
var markAsRead bool

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
		if err := core.ClearScreen(); err != nil {
			slog.Error("Failed to clear screen", "error", err)
			os.Exit(1)
		}

		// ------------ get entries ------------ //
		entries, err := core.GetEntries()
		if err != nil {
			slog.Error("Cannot obtain articles from Wallabag")
			os.Exit(1)
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
			slog.Error("Form error", "error", err)
			os.Exit(1)
		}

		// ------------ summarize ------------ //
		fmt.Printf("========== %s ==========\n", entryTitle)

		var content string
		var entryID int
		for _, entry := range entries {
			if entry.Title == entryTitle {
				content = entry.Content
				entryID = entry.ID
			}
		}

		p := bluemonday.StripTagsPolicy()
		contentSanitized := p.Sanitize(
			content,
		)

		if _, err := core.Summarize(contentSanitized, core.DetectLanguage(content), "cli"); err != nil {
			slog.Error("Failed to summarize article", "error", err)
			os.Exit(1)
		}

		formMarkAsRead := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Mark article as read?").
					Value(&markAsRead).
					Affirmative("Yes").
					Negative("No"),
			),
		)
		if err := formMarkAsRead.Run(); err != nil {
			slog.Error("Form error", "error", err)
			os.Exit(1)
		}

		if markAsRead {
			if err := core.MarkEntryAsRead(entryID); err != nil {
				slog.Error("Failed to mark entry as read", "error", err)
				os.Exit(1)
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("Command failed", "error", err)
		os.Exit(1)
	}
}

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	logger := zerolog.New(output).With().Timestamp().Logger()
	slog.SetDefault(slog.New(slogzerolog.Option{Logger: &logger}.NewZerologHandler()))

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
