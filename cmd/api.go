/*
Copyright © 2024 Karn Wong <karn@karnwong.me>
*/
package cmd

import (
	"net/http"
	"os"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/kahnwong/article-summarizer/core"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"
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
		log := zerolog.New(os.Stderr).With().Timestamp().Logger()

		app.Use(logger.SetLogger(logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return log
		})))

		// routes
		app.GET("/", rootController)

		// error handling
		if err := app.Run(":3000"); err != nil {
			log.Fatal().Msg("Gin app error")
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
