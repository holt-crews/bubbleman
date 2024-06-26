package requestbody

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textarea.Model
}

func (m *Model) View() string {
	return m.Model.View()
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

type Option func(*Model)

func New(opts ...Option) *Model {
	t := textarea.New()
	t.ShowLineNumbers = true
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()

	model := &Model{
		Model: t,
	}

	for _, opt := range opts {
		opt(model)
	}

	return model
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

func WithFocusedStyle(style textarea.Style) Option {
	return func(m *Model) {
		m.FocusedStyle = style
	}
}

func WithBlurredStyle(style textarea.Style) Option {
	return func(m *Model) {
		m.BlurredStyle = style
	}
}
