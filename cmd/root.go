/*
Copyright © 2026 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/Strubbl/wallabago/v9"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"github.com/spf13/cobra"

	"charm.land/huh/v2"
	"github.com/kahnwong/article-summarizer/core"
)

var entryTitle string
var markAsRead bool

var wallabagClient core.WallabagClient

// functions
func createFormOptions(entries []wallabago.Item) []huh.Option[string] {
	var options []huh.Option[string]
	gray := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	for _, v := range entries {
		articleURL, err := url.Parse(v.URL)
		if err != nil || articleURL.Scheme == "" || articleURL.Host == "" {
			options = append(options, huh.NewOption(v.Title, v.Title))
			continue
		}

		label := fmt.Sprintf("%s %s", v.Title, gray.Render("("+articleURL.Host+")"))
		options = append(options, huh.NewOption(label, v.Title))
	}

	return options
}

// main
var rootCmd = &cobra.Command{
	Use:   "article-summarizer",
	Short: "Summarize an article with LLM",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			// Clears the screen
			if err := core.ClearScreen(); err != nil {
				slog.Error("Failed to clear screen", "error", err)
				os.Exit(1)
			}

			// ------------ get entries ------------ //
			entries, err := wallabagClient.GetEntries()
			if err != nil {
				slog.Error("Cannot obtain articles from Wallabag")
				os.Exit(1)
			}

			// ------------ select article ------------ //
			entryTitle = ""
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
				if errors.Is(err, huh.ErrUserAborted) {
					return
				}
				slog.Error("Form error", "error", err)
				os.Exit(1)
			}

			// ------------ summarize ------------ //
			fmt.Printf("========== %s ==========\n", entryTitle)

			var selectedEntry wallabago.Item
			for _, entry := range entries {
				if entry.Title == entryTitle {
					selectedEntry = entry
				}
			}

			if _, err := core.SummarizeArticle(selectedEntry, "cli"); err != nil {
				slog.Error("Failed to summarize article", "error", err)
				os.Exit(1)
			}

			// ------------ mark as read ------------ //
			markAsRead = false
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
				if errors.Is(err, huh.ErrUserAborted) {
					return
				}
				slog.Error("Form error", "error", err)
				os.Exit(1)
			}

			if !markAsRead {
				return
			}

			if err := wallabagClient.MarkEntryAsRead(selectedEntry.ID); err != nil {
				slog.Error("Failed to mark entry as read", "error", err)
				os.Exit(1)
			}
		}
	},
}

func Execute() {
	if err := configureLogger(); err != nil {
		slog.Error("Failed to configure logger", "error", err)
		os.Exit(1)
	}

	if err := core.LoadConfig(); err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	wallabagClient = core.NewWallabagClient()

	err := rootCmd.Execute()
	if err != nil {
		slog.Error("Command failed", "error", err)
		os.Exit(1)
	}
}

func configureLogger() error {
	level, err := parseLogLevel(os.Getenv("LOG_LEVEL"))

	output := zerolog.ConsoleWriter{Out: os.Stderr}
	zerologLogger := zerolog.New(output)
	logger := slog.New(slogzerolog.Option{
		Level:  level,
		Logger: &zerologLogger,
	}.NewZerologHandler())
	slog.SetDefault(logger)

	return err
}

func parseLogLevel(value string) (slog.Level, error) {
	level := slog.LevelDebug
	if value == "" {
		return level, nil
	}

	if err := level.UnmarshalText([]byte(value)); err != nil {
		return slog.LevelDebug, fmt.Errorf("invalid LOG_LEVEL %q: %w", value, err)
	}
	return level, nil
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
