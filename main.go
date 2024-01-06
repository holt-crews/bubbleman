package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	initialInputs = 2
	helpHeight    = 5
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

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
	t.BlurredStyle.Base = blurredBorderStyle
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
	m.requestBody.SetWidth(m.width)
	return m
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.request):
			cmd := m.requestBody.Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.url):
			m.requestBody.Blur()
			cmd := m.urlbar.Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.send):
			m.sendRequest()
			// cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}
	m.sizeInputs()

	// update text area
	newRequest, cmd := m.requestBody.Update(msg)
	m.requestBody = newRequest
	cmds = append(cmds, cmd)

	newUrl, cmd := m.urlbar.Update(msg)
	m.urlbar = newUrl
	cmds = append(cmds, cmd)

	// Update all text
	return m, tea.Batch(cmds...)
}

func (m *model) sizeInputs() {
	m.requestBody.SetWidth(m.width - 2)
	m.requestBody.SetHeight(m.height - helpHeight)
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
	)
	doc.WriteString(requestInputs)
	doc.WriteString("\n")
	doc.WriteString(help)

	return docStyle.Render(doc.String())
}

func (m model) sendRequest() *http.Response {
	resp, err := http.Get(m.urlbar.Value())
	if err != nil {
		log.Fatalln(err)
	}

	// just doing this for now to prove that it's working, look into logging in the future
	fmt.Println("response: ", &resp.Body)
	os.Exit(1)

	return resp
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
