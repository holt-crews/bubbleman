package urlbar

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textinput.Model
}

func (u *Model) View() string {
	return u.Model.View()
}

func (u *Model) Update(msg tea.Msg) tea.Cmd {
	t, cmd := u.Model.Update(msg)
	u.Model = t
	return cmd
}

func (u *Model) Value() string {
	return u.Model.Value()
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
	return func(u *Model) {
		u.Prompt = prompt
	}
}

func WithPlaceholder(placeholder string) Option {
	return func(u *Model) {
		u.Placeholder = placeholder
	}
}

func WithCursorStyle(style lipgloss.Style) Option {
	return func(u *Model) {
		u.Cursor.Style = style
	}
}
