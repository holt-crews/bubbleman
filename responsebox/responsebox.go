package responsebox

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResponseBox struct {
	viewport.Model
}

func (u *ResponseBox) View() string {
	return u.Model.View()
}

func (u *ResponseBox) Update(msg tea.Msg) tea.Cmd {
	t, cmd := u.Model.Update(msg)
	u.Model = t
	return cmd
}

// TODO
func (u *ResponseBox) Value() string {
	return ""
}

type ResponseBoxOption func(*ResponseBox)

func NewResponseBox(opts ...ResponseBoxOption) *ResponseBox {
	t := viewport.New(1, 1)
	t.MouseWheelEnabled = true

	responsebox := &ResponseBox{
		Model: t,
	}

	for _, opt := range opts {
		opt(responsebox)
	}

	return responsebox
}

func WithPlaceholder(placeholder string) ResponseBoxOption {
	return func(u *ResponseBox) {
		u.SetContent(placeholder)
	}
}

func WithStyle(style lipgloss.Style) ResponseBoxOption {
	return func(u *ResponseBox) {
		u.Style = style
	}
}
