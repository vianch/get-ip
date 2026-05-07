package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Refresh key.Binding
	WiFi    key.Binding
	Quit    key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		WiFi: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "wifi"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Refresh, k.WiFi, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Refresh, k.WiFi, k.Quit}}
}
