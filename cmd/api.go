/*
Copyright © 2026 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"log/slog"
	"net/http"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/kahnwong/article-summarizer/core"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/cobra"
)

func rootController(c *gin.Context) {
	// ------------ get entries ------------ //
	entries, err := core.GetEntries()
	if err != nil {
		c.String(http.StatusInternalServerError, "cannot obtain articles from Wallabag: %v", err)
		return
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

	output, err := core.Summarize(contentSanitized, core.DetectLanguage(content), "api")
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to summarize article: %v", err)
		return
	}

	c.String(http.StatusOK, "===== %s =====\n%s", title, output)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Serve summarization as API",
	Run: func(cmd *cobra.Command, args []string) {
		// app
		app := gin.New()
		app.Use(logger.SetLogger())

		// routes
		app.GET("/", rootController)

		// error handling
		if err := app.Run(":3000"); err != nil {
			slog.Error("Gin app error", "error", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
