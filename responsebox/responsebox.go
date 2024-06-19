package responsebox

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport.Model
	FocusedStyle lipgloss.Style
	BlurredStyle lipgloss.Style

	focus bool
}

func (m *Model) View() string {
	return m.Model.View()
}

func (m *Model) Blur() {
	m.focus = false
	m.Style = m.BlurredStyle
}

func (m *Model) Focus() tea.Cmd {
	m.focus = true
	m.Style = m.FocusedStyle
	return nil
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	t, cmd := m.Model.Update(msg)
	m.Model = t
	return cmd
}

// TODO: implement
func (m *Model) Value() string {
	return ""
}

func (m *Model) SetWidth(w int) {
	m.Width = w
}

func (m *Model) SetHeight(h int) {
	m.Height = h
}

type Option func(*Model)

func New(opts ...Option) *Model {
	t := viewport.New(1, 1)
	t.MouseWheelEnabled = true

	responsebox := &Model{
		Model: t,
	}

	for _, opt := range opts {
		opt(responsebox)
	}

	return responsebox
}

func WithPlaceholder(placeholder string) Option {
	return func(m *Model) {
		m.SetContent(placeholder)
	}
}

func WithStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.Style = style
	}
}

func WithFocusedStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.FocusedStyle = style
	}
}

func WithBlurredStyle(style lipgloss.Style) Option {
	return func(m *Model) {
		m.BlurredStyle = style
	}
}
