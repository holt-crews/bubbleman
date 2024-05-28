package responsebox

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport.Model
}

func (u *Model) View() string {
	return u.Model.View()
}

func (u *Model) Update(msg tea.Msg) tea.Cmd {
	t, cmd := u.Model.Update(msg)
	u.Model = t
	return cmd
}

// TODO
func (u *Model) Value() string {
	return ""
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
	return func(u *Model) {
		u.SetContent(placeholder)
	}
}

func WithStyle(style lipgloss.Style) Option {
	return func(u *Model) {
		u.Style = style
	}
}
