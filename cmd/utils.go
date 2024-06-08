package cmd

import "golang.design/x/clipboard"

func ReadFromClipboard() string {
	err := clipboard.Init()
	if err == nil {

		s := clipboard.Read(clipboard.FmtText)
		return string(s)
	}

	return ""
}
