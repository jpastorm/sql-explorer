package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func extractCurrentLine(m textarea.Model) string {
	lines := strings.Split(m.Value(), "\n")
	cursorLine := m.Line()

	if cursorLine >= len(lines) {
		return ""
	}

	return strings.TrimSpace(lines[cursorLine])
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
