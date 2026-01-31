// Package app: here placed root model
package app

import (
	nav "client/internal/navigator"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	nav *nav.Navigator
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func New(nav *nav.Navigator) *Model {
	return &Model{
		nav: nav,
	}
}
