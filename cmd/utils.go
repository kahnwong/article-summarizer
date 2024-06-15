package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

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
