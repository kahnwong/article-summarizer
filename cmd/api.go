/*
Copyright © 2026 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/kahnwong/article-summarizer/core"
	slogfiber "github.com/samber/slog-fiber"
	"github.com/spf13/cobra"
)

func rootController(c fiber.Ctx) error {
	// ------------ get entries ------------ //
	entries, err := wallabagClient.GetEntries()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf("cannot obtain articles from Wallabag: %v", err),
		)
	}

	// ------------ get title and content ------------ //
	entry := entries[0]

	// ------------ summarize ------------ //
	output, err := core.SummarizeArticle(entry, "api")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf("failed to summarize article: %v", err),
		)
	}

	return c.SendString(fmt.Sprintf("===== %s =====\n%s", entry.Title, output))
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Serve summarization as API",
	Run: func(cmd *cobra.Command, args []string) {
		// app
		app := fiber.New()
		app.Use(slogfiber.New(slog.Default()))
		app.Use(recover.New())

		// routes
		app.Get("/", rootController)

		// error handling
		if err := app.Listen(":3000"); err != nil {
			slog.Error("Fiber app error", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
