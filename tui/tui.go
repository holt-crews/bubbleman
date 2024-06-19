package tui

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/selection"

	"github.com/holt-crews/bubbleman/helpers"
	"github.com/holt-crews/bubbleman/requestbody"
	"github.com/holt-crews/bubbleman/responsebox"
	"github.com/holt-crews/bubbleman/urlbar"
)

type Adjustable interface {
	SetWidth(w int)
	SetHeight(w int)
}

type Selectable interface {
	Focus() tea.Cmd
	Blur()
}

type Component interface {
	// tea.Model
	// Init() tea.Cmd
	View() string
	Update(msg tea.Msg) tea.Cmd
	Value() string

	Selectable
	Adjustable
	// SetDimensions()
	// Select() // used when "enter" is hit on that component
}

type model struct {
	help         help.Model
	requestBody  Component
	urlbar       Component
	response     Component
	selection    *selection.Model[string]
	httpMethod   string
	keymap       keymap
	components   []Component
	viewReady    bool
	methodToggle bool
}

func initialModel() model {
	sel := selection.New("Items", httpMethods)
	sel.Filter = nil

	urlbar := urlbar.New(
		urlbar.WithPrompt(""),
		urlbar.WithPlaceholder("https://api.com/v1"),
	)
	requestBody := requestbody.New(
		requestbody.WithPrompt(""),
		requestbody.WithPlaceholder("Type something..."),
		requestbody.WithCursorStyle(cursorStyle),
		requestbody.WithFocusedStyle(focusedTextAreaStyle),
		requestbody.WithBlurredStyle(blurredTextAreaStyle),
	)

	m := model{
		urlbar:      urlbar,
		httpMethod:  "GET",
		requestBody: requestBody,
		help:        help.New(),
		keymap:      Keymap,
		selection:   selection.NewModel(sel),
		components:  []Component{urlbar, requestBody},
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

	rCmd = m.requestBody.Update(msg)
	uCmd = m.urlbar.Update(msg)
	respCmd = m.response.Update(msg)
	cmds = append(cmds, rCmd, uCmd, respCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Request):
			// m.urlbar.Blur()
			cmd := m.requestBody.Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.Url):
			m.requestBody.Blur()
			// cmd := m.urlbar.Focus()
			// cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.HttpMethod):
			m.methodToggle = !m.methodToggle
		case key.Matches(msg, m.keymap.Send):
			// m.urlbar.Blur()
			m.requestBody.Blur()
			// TODO: fix
			// resp := m.sendRequest()
			// m.response.SetContent(resp)
			// cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		WindowSize = msg
		if !m.viewReady {
			m.response = responsebox.New(
				responsebox.WithPlaceholder("{\"response\": \"OK\"...}"),
				responsebox.WithStyle(focusedBorderStyle.Foreground(gray)),
			)
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

	m.urlbar.SetWidth(remainingWidth)

	m.requestBody.SetWidth(remainingWidth)
	// gotta be careful with division because it floors it, the divide by 3 messes with it, consider alternate ways
	m.requestBody.SetHeight(2*(remainingHeight/3) - 2)
	// there's a bug in viewport: https://github.com/charmbracelet/bubbles/pull/388
	// TODO: ideally the "- 2" is dynamic based on the height of the http method box

	// TODO: fix
	m.response.SetHeight((remainingHeight / 3) - 2)
	// .SetWidth() and .Width are calculated differently. 2 seems to be magic difference for my case
	m.response.SetWidth(remainingWidth + 2)
}

func (m model) View() string {
	doc := strings.Builder{}
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.Request,
		m.keymap.Url,
		m.keymap.Send,
		m.keymap.Quit,
	})

	// TODO: think about if there's a way to place this style information on the component
	// and then we could have a function that we just call for each element in the grid
	// right now it's kinda tricky because "requestInputs" is dependent on "bar"
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
