package tui

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/selection"

	"github.com/holt-crews/bubbleman/helpers"
)

// consider just wrapping everything into a bubble tea Model in a "components package"
type Component interface {
	// Init() tea.Cmd
	// for some unknown reason, bubbles don't actually implement the tea.Model interface, missing Init()
	// Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
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
	// currently faded text but will want color when a response is received
	m.response.Style = focusedBorderStyle.Foreground(gray)
	m.response.SetContent(fmt.Sprintf("{\"response\": \"OK\"...}"))
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
	keymap       keymap
	help         help.Model
	requestBody  textarea.Model
	httpMethod   string
	urlbar       textinput.Model
	response     viewport.Model
	viewReady    bool
	selection    *selection.Model[string]
	methodToggle bool
	components   []Component
	// focus       int  // will probably want to come back to this when all components are laid out
}

func InitialRequest() model {
	sel := selection.New("Items", httpMethods)
	sel.Filter = nil

	m := model{
		urlbar:      newUrlbar(),
		httpMethod:  "GET",
		requestBody: newTextarea(),
		help:        help.New(),
		keymap:      Keymap,
		selection:   selection.NewModel(sel),
		components:  []Component{newUrlbar(), newTextarea()},
	}
	m.selection.Init()

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
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Request):
			m.urlbar.Blur()
			cmd := m.requestBody.Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.Url):
			m.requestBody.Blur()
			cmd := m.urlbar.Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.HttpMethod):
			m.methodToggle = !m.methodToggle
		case key.Matches(msg, m.keymap.Send):
			m.urlbar.Blur()
			m.requestBody.Blur()
			resp := m.sendRequest()
			m.response.SetContent(resp)
			// cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		WindowSize = msg
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
	top, right, bottom, left := DocStyle.GetMargin()
	remainingHeight := WindowSize.Height - top - bottom
	remainingWidth := WindowSize.Width - left - right

	m.urlbar.Width = remainingWidth

	m.requestBody.SetWidth(remainingWidth)
	// gotta be careful with division because it floors it, the divide by 3 messes with it, consider alternate ways
	m.requestBody.SetHeight(2*(remainingHeight/3) - 2)
	// there's a bug in viewport: https://github.com/charmbracelet/bubbles/pull/388
	// TODO: ideally the "- 2" is dynamic based on the height of the http method box
	m.response.Height = (remainingHeight / 3) - 2
	// .SetWidth() and .Width are calculated differently. 2 seems to be magic difference for my case
	m.response.Width = remainingWidth + 2
}

func (m model) View() string {
	doc := strings.Builder{}
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.Request,
		m.keymap.Url,
		m.keymap.Send,
		m.keymap.Quit,
	})

	bar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		methodBoxStyle.Render(m.httpMethod),
		m.urlbar.View(),
	)
	lipgloss.Height(bar)
	requestInputs := lipgloss.JoinVertical(
		lipgloss.Top,
		bar,
		m.requestBody.View(),
		m.response.View(),
	)
	doc.WriteString(requestInputs)
	if m.methodToggle {
		overlay := helpers.PlaceOverlay(10, 10, m.selection.View(), doc.String(), false)
		doc.WriteString(overlay)
	}
	doc.WriteString("\n")
	doc.WriteString(help)

	return DocStyle.Render(doc.String())
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
