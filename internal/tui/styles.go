package tui

import "github.com/charmbracelet/lipgloss"

var (
	accent = lipgloss.Color("#7D56F4")
	dim    = lipgloss.Color("241")
	red    = lipgloss.Color("203")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(1, 2).
			Foreground(accent)

	tableBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(accent)

	// Padding(0, 1) matches bubbles/table's default Cell padding so header
	// text aligns with row text — without it, every header sits one column
	// to the left of its data and the rightmost cell's bottom border drops.
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(accent)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(accent).
			Bold(true)

	helpStyle = lipgloss.NewStyle().Foreground(dim)

	wifiPanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Padding(0, 2).
			MarginTop(1)

	wifiLabelStyle = lipgloss.NewStyle().Foreground(accent).Bold(true)

	statusErrStyle = lipgloss.NewStyle().Foreground(red)
)
