package core

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func ClearScreen() error {
	// from <https://github.com/MasterDimmy/go-cls/blob/main/cls.go>
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear") //Linux example, its tested
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls") //Windows example, its tested
	default:
		return fmt.Errorf("CLS for %s not implemented", runtime.GOOS)
	}

	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clear screen: %w", err)
	}

	return nil
}
