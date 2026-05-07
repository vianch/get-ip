package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/vianch/get-ip/internal/tui"
)

var version = "dev"

func main() {
	for _, a := range os.Args[1:] {
		if a == "-v" || a == "--version" {
			fmt.Println("get-ip", version)
			return
		}
		if a == "-h" || a == "--help" {
			fmt.Println("get-ip — TUI for local network interfaces")
			fmt.Println("Keys: r=refresh w=wifi q=quit")
			return
		}
	}

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "get-ip:", err)
		os.Exit(1)
	}
}
