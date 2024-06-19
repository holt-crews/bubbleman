package urlbar

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textinput.Model
}

func (m *Model) View() string {
	return m.Model.View()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	t, cmd := m.Model.Update(msg)
	m.Model = t
	return cmd
}

func (m *Model) Value() string {
	return m.Model.Value()
}

func (m *Model) SetWidth(w int) {
	m.Width = w
}

func (m *Model) SetHeight(h int) {
	// can't set height
}

type Option func(*Model)

func New(opts ...Option) *Model {
	t := textinput.New()
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.Blur()

	urlbar := &Model{
		Model: t,
	}

	for _, opt := range opts {
		opt(urlbar)
	}

	return urlbar
}

func WithPrompt(prompt string) Option {
	return func(m *Model) {
		m.Prompt = prompt
	}
}

func WithPlaceholder(placeholder string) Option {
	return func(m *Model) {
		m.Placeholder = placeholder
	}
}

func WithCursorStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.Cursor.Style = style
	}
}
