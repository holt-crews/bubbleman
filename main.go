package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	initialInputs = 2
	helpHeight    = 5
	rightPadding  = 4

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

type keymap = struct {
	request, url, send, quit key.Binding
}

func newUrlbar() textinput.Model {
	t := textinput.New()
	t.Prompt = ""
	t.Placeholder = "https://api.something.com/v1/users"
	t.Cursor.Style = cursorStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.Blur()
	return t
}

func (m *model) newResponseView() viewport.Model {
	m.response = viewport.New(1, 1)
	m.response.Style = focusedBorderStyle
	m.response.SetContent(fmt.Sprintf("width: %d; height: %d", m.width, m.height))
	m.response.MouseWheelEnabled = true
	return m.response
}

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = focusedBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}

type model struct {
	width       int
	height      int
	keymap      keymap
	help        help.Model
	requestBody textarea.Model
	urlbar      textinput.Model
	response    viewport.Model
	viewReady   bool
	// focus       int  // will probably want to come back to this when all components are laid out
}

func initialModel() model {
	m := model{
		urlbar:      newUrlbar(),
		requestBody: newTextarea(),
		help:        help.New(),
		keymap: keymap{
			request: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "edit request body"),
			),
			url: key.NewBinding(
				key.WithKeys("ctrl+u"),
				key.WithHelp("ctrl+u", "edit url bar"),
			),
			send: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "send request"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}

	// initially focus on requestBody section
	m.requestBody.Focus()
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		rCmd    tea.Cmd
		uCmd    tea.Cmd
		respCmd tea.Cmd
	)
	var cmds []tea.Cmd

	m.requestBody, rCmd = m.requestBody.Update(msg)
	m.urlbar, uCmd = m.urlbar.Update(msg)
	m.response, respCmd = m.response.Update(msg)
	cmds = append(cmds, rCmd, uCmd, respCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.request):
			m.urlbar.Blur()
			cmd := m.requestBody.Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.url):
			m.requestBody.Blur()
			cmd := m.urlbar.Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.send):
			m.urlbar.Blur()
			m.requestBody.Blur()
			resp := m.sendRequest()
			m.response.SetContent(resp)
			// cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		if !m.viewReady {
			m.newResponseView()
			m.viewReady = true
		}
		m.sizeInputs()
	}

	// Update all text
	return m, tea.Batch(cmds...)
}

// needs to be a pointer receiver in order to update
func (m *model) sizeInputs() {
	totalHeight := m.height - (lipgloss.Height(m.urlbar.View()) + helpHeight)

	m.urlbar.Width = m.width - rightPadding

	m.requestBody.SetWidth(m.width - rightPadding)
	m.requestBody.SetHeight(2 * (totalHeight / 3))

	m.response.Height = totalHeight / 3
	m.response.Width = m.width - rightPadding
}

func (m model) View() string {
	doc := strings.Builder{}
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.request,
		m.keymap.url,
		m.keymap.send,
		m.keymap.quit,
	})

	requestInputs := lipgloss.JoinVertical(
		lipgloss.Top,
		"GET "+m.urlbar.View(),
		m.requestBody.View(),
		m.response.View(),
	)
	doc.WriteString(requestInputs)
	doc.WriteString("\n")
	doc.WriteString(help)

	return docStyle.Render(doc.String())
}

func (m model) sendRequest() string {
	resp, err := http.Get(m.urlbar.Value())
	if err != nil {
		log.Fatalln(err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	return bodyString
}

func main() {
	// WithMouse is not working in viewport for whatever reason
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
