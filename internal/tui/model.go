package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/vianch/get-ip/internal/network"
)

type refreshedMsg struct {
	rows []network.Interface
	err  error
}

type wifiMsg struct {
	info network.WiFiInfo
	err  error
}

type Model struct {
	table       table.Model
	help        help.Model
	keys        keyMap
	rows        []network.Interface
	err         error
	loading     bool
	wifiOpen    bool
	wifi        *network.WiFiInfo
	wifiLoading bool
	wifiErr     error
}

func New() Model {
	cols := []table.Column{
		{Title: "Device", Width: 22},
		{Title: "IPv4", Width: 16},
		{Title: "MAC", Width: 19},
		{Title: "Gateway", Width: 16},
		{Title: "Mode", Width: 8},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(8),
	)
	st := table.DefaultStyles()
	st.Header = headerStyle
	st.Selected = selectedStyle
	t.SetStyles(st)

	h := help.New()
	h.Styles.ShortKey = helpStyle
	h.Styles.ShortDesc = helpStyle
	h.Styles.ShortSeparator = helpStyle

	return Model{
		table:   t,
		help:    h,
		keys:    newKeyMap(),
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return refreshCmd
}

func refreshCmd() tea.Msg {
	rows, err := network.Collect()
	return refreshedMsg{rows: rows, err: err}
}

func wifiCmd() tea.Msg {
	info, err := network.WiFi()
	return wifiMsg{info: info, err: err}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.wifi = nil
			cmds := []tea.Cmd{refreshCmd}
			if m.wifiOpen {
				m.wifiLoading = true
				cmds = append(cmds, wifiCmd)
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.WiFi):
			m.wifiOpen = !m.wifiOpen
			if m.wifiOpen && m.wifi == nil && !m.wifiLoading {
				m.wifiLoading = true
				return m, wifiCmd
			}
			return m, nil
		}

	case refreshedMsg:
		m.loading = false
		m.err = msg.err
		m.rows = msg.rows
		m.table.SetRows(toTableRows(msg.rows))
		return m, nil

	case wifiMsg:
		m.wifiLoading = false
		if msg.err != nil {
			m.wifiErr = msg.err
			return m, nil
		}
		info := msg.info
		m.wifi = &info
		m.wifiErr = nil
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func toTableRows(in []network.Interface) []table.Row {
	rows := make([]table.Row, 0, len(in))
	for _, i := range in {
		mac := i.MAC
		if mac == "" {
			mac = "—"
		}
		gw := i.Gateway
		if gw == "" {
			gw = "—"
		}
		// bubbles/table measures cell width by raw byte length, so embedding
		// ANSI escapes here would eat the column budget and break alignment.
		rows = append(rows, table.Row{i.Label(), i.IPv4, mac, gw, string(i.Mode)})
	}
	return rows
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🌐 get-ip"))
	b.WriteString("\n")

	if m.loading && len(m.rows) == 0 {
		b.WriteString(helpStyle.Render("  Scanning interfaces…"))
		b.WriteString("\n\n")
	} else {
		b.WriteString(tableBorderStyle.Render(m.table.View()))
		b.WriteString("\n")
	}

	if m.err != nil {
		b.WriteString(statusErrStyle.Render("  error: " + m.err.Error()))
		b.WriteString("\n")
	}

	if m.wifiOpen {
		b.WriteString(m.wifiView())
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  " + m.help.ShortHelpView(m.keys.ShortHelp())))
	b.WriteString("\n")

	return b.String()
}

func (m Model) wifiView() string {
	if m.wifiLoading {
		return wifiPanelStyle.Render(helpStyle.Render("Loading WiFi…"))
	}
	if m.wifiErr != nil {
		return wifiPanelStyle.Render(statusErrStyle.Render("WiFi: " + m.wifiErr.Error()))
	}
	if m.wifi == nil {
		return ""
	}
	w := m.wifi
	label := wifiLabelStyle.Render("WiFi")
	parts := []string{
		label,
		fmt.Sprintf("SSID: %s", w.SSID),
	}
	if w.Channel != "" {
		parts = append(parts, "Channel: "+w.Channel)
	}
	if w.Signal != "" {
		parts = append(parts, "Signal: "+w.Signal)
	}
	line := strings.Join(parts, lipgloss.NewStyle().Foreground(dim).Render(" • "))
	return wifiPanelStyle.Render(line)
}
