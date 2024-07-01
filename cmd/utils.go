package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.design/x/clipboard"
)

func ReadFromClipboard() string {
	err := clipboard.Init()
	if err == nil {

		s := clipboard.Read(clipboard.FmtText)
		return string(s)
	}

	return ""
}

func ClearScreen() {
	// from <https://github.com/MasterDimmy/go-cls/blob/main/cls.go>

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear") //Linux example, its tested
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") //Windows example, its tested
	default:
		fmt.Println("CLS for ", runtime.GOOS, " not implemented")
		return
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func getUrl(args []string) string {
	var url string

	if len(args) == 0 {
		urlFromClipboard := ReadFromClipboard()
		if urlFromClipboard != "" {
			if strings.HasPrefix(urlFromClipboard, "https://") {
				url = urlFromClipboard
			}
		}
	}
	if url == "" {
		if len(args) == 0 {
			fmt.Println("Please specify URL")
			os.Exit(1)
		} else if len(args) == 1 {
			url = args[0]
		}
	}

	return url
}
