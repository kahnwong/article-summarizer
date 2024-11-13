/*
Copyright © 2024 Karn Wong <karn@karnwong.me>
*/
package main

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/kahnwong/article-summarizer/cmd"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
)

// import "github.com/kahnwong/article-summarizer/cmd"
//
//	func main() {
//		cmd.Execute()
//	}
func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c fiber.Ctx) error {
		// ------------ get entries ------------ //
		entries, err := cmd.GetEntries()
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot obtain articles from Wallabag")
		}
		entry := entries[0]

		// ------------ summarize ------------ //
		entryTitle := entry.Title
		content := entry.Content

		p := bluemonday.StripTagsPolicy()
		contentSanitized := p.Sanitize(
			content,
		)

		response := cmd.Summarize(contentSanitized, cmd.DetectLanguage(content))

		return c.SendString(fmt.Sprintf("%s\n=====\n%s", entryTitle, response))
	})

	// Start the server on port 3000
	log.Fatal().Err(app.Listen(":3000"))
}
