/*
Copyright © 2024 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/kahnwong/article-summarizer/core"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func rootController(c *fiber.Ctx) error {
	// ------------ get entries ------------ //
	entries, err := core.GetEntries()
	if err != nil {
		log.Fatal().Msg("Cannot obtain articles from Wallabag")
	}

	// ------------ get title and content ------------ //
	entry := entries[0]
	title := entry.Title
	content := entry.Content

	// ------------ summarize ------------ //
	p := bluemonday.StripTagsPolicy()
	contentSanitized := p.Sanitize(
		content,
	)

	output := core.Summarize(contentSanitized, core.DetectLanguage(content), "api")

	return c.SendString(fmt.Sprintf("===== %s =====\n%s", title, output))
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Serve summarization as API",
	Run: func(cmd *cobra.Command, args []string) {
		// app
		app := fiber.New()
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

		app.Use(fiberzerolog.New(fiberzerolog.Config{
			Logger: &logger,
		}))

		// routes
		app.Get("/", rootController)

		// error handling
		if err := app.Listen(":3000"); err != nil {
			logger.Fatal().Msg("Fiber app error")
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
