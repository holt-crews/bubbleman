package tui

import "github.com/charmbracelet/lipgloss"

const (
	helpHeight   = 2
	rightPadding = 6

	// TODO: pull in from config file
	black   = "#928374"
	blue    = "#83a598"
	cyan    = "#8ec07c"
	green   = "#b8bb26"
	magenta = "#d3869b"
	red     = "#fb4934"
	white   = "#ebdbb2"
	yellow  = "#fabd2f"
)

var (
	httpMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(white))

	cursorLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(cyan))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(cyan))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(cyan))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(cyan))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(blue))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)
