package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UrlBar struct {
	textinput.Model
}

func (u *UrlBar) View() string {
	return u.Model.View()
}

func (u *UrlBar) Update(msg tea.Msg) tea.Cmd {
	t, cmd := u.Model.Update(msg)
	u.Model = t
	return cmd
}

func (u *UrlBar) Value() string {
	return u.Model.Value()
}

type UrlBarOption func(*UrlBar)

func NewUrlBar(opts ...UrlBarOption) *UrlBar {
	t := textinput.New()
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.Blur()

	urlbar := &UrlBar{
		Model: t,
	}

	for _, opt := range opts {
		opt(urlbar)
	}

	return urlbar
}

func WithPrompt(prompt string) UrlBarOption {
	return func(u *UrlBar) {
		u.Prompt = prompt
	}
}

func WithPlaceholder(placeholder string) UrlBarOption {
	return func(u *UrlBar) {
		u.Placeholder = placeholder
	}
}

func WithCursorStyle(style lipgloss.Style) UrlBarOption {
	return func(u *UrlBar) {
		u.Cursor.Style = style
	}
}
