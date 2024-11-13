/*
Copyright © 2024 Karn Wong <karn@karnwong.me>
*/
package cmd

var AppConfig = readConfig()

// var entryTitle string // for huh form
//
// var rootCmd = &cobra.Command{
// 	Use:   "article-summarizer",
// 	Short: "Summarize an article with LLM",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// Clears the screen
// 		ClearScreen()
//
// 		// ------------ get entries ------------ //
// 		entries, err := GetEntries()
// 		if err != nil {
// 			log.Fatal().Err(err).Msg("Cannot obtain articles from Wallabag")
// 		}
//
// 		// ------------ select article ------------ //
// 		formEntries := huh.NewForm(
// 			huh.NewGroup(
// 				huh.NewSelect[string]().
// 					Title("Choose an article to summarize").
// 					Options(
// 						createFormOptions(entries)...,
// 					).
// 					Value(&entryTitle),
// 			),
// 		)
// 		err = formEntries.Run()
// 		if err != nil {
// 			log.Fatal().Err(err)
// 		}
//
// 		// ------------ summarize ------------ //
// 		fmt.Printf("========== %s ==========\n", entryTitle)
//
// 		var content string
// 		for _, entry := range entries {
// 			if entry.Title == entryTitle {
// 				content = entry.Content
// 			}
// 		}
//
// 		p := bluemonday.StripTagsPolicy()
// 		contentSanitized := p.Sanitize(
// 			content,
// 		)
//
// 		Summarize(contentSanitized, DetectLanguage(content))
// 	},
// }
//
// func Execute() {
// 	err := rootCmd.Execute()
// 	if err != nil {
// 		os.Exit(1)
// 	}
// }
//
// func init() {
// 	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
