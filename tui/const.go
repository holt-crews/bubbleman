package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// TODO: make this dynamic somehow
	helpHeight = 2

	// TODO: pull in from config file
	black   = lipgloss.Color("#928374")
	blue    = lipgloss.Color("#83a598")
	cyan    = lipgloss.Color("#8ec07c")
	green   = lipgloss.Color("#b8bb26")
	magenta = lipgloss.Color("#d3869b")
	red     = lipgloss.Color("#fb4934")
	white   = lipgloss.Color("#ebdbb2")
	yellow  = lipgloss.Color("#fabd2f")
	gray    = lipgloss.Color("240") // mostly for placeholder text
)

var (
	httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}

	cursorStyle             = lipgloss.NewStyle().Foreground(white)
	cursorLineStyle         = lipgloss.NewStyle().Foreground(cyan)
	placeholderStyle        = lipgloss.NewStyle().Foreground(cyan)
	endOfBufferStyle        = lipgloss.NewStyle().Foreground(cyan)
	focusedPlaceholderStyle = lipgloss.NewStyle().Foreground(cyan)

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(blue)
	blurredBorderStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
	methodBoxStyle     = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(blue).Padding(0, 1).Margin(0, 1, 0, 0)

	DocStyle = lipgloss.NewStyle().Margin(2, 3, 0, 1)

	WindowSize tea.WindowSizeMsg

	focusedTextAreaStyle = textarea.Style{
		Placeholder: focusedPlaceholderStyle,
		CursorLine:  cursorLineStyle,
		Base:        focusedBorderStyle,
		EndOfBuffer: endOfBufferStyle,
	}
	blurredTextAreaStyle = textarea.Style{
		Placeholder: placeholderStyle,
		Base:        focusedBorderStyle,
		EndOfBuffer: endOfBufferStyle,
	}
)

type keymap = struct {
	Request, Url, HttpMethod, Send, Quit key.Binding
}

var Keymap = keymap{
	Request: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "edit request body"),
	),
	Url: key.NewBinding(
		key.WithKeys("ctrl+b"),
		key.WithHelp("ctrl+b", "edit url bar"),
	),
	HttpMethod: key.NewBinding(
		key.WithKeys("ctrl+u"),
		key.WithHelp("ctrl+u", "edit http method"),
	),
	Send: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "send request"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
