// Package app: here placed root model
package app

import (
	nav "client/internal/navigator"

	tea "github.com/charmbracelet/bubbletea"
)

type UserState struct {
	// todo: add here auth token logic
}

type Model struct {
	nav   *nav.Navigator
	state *UserState
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func New(nav *nav.Navigator, state *UserState) *Model {
	return &Model{
		nav:   nav,
		state: state,
	}
}
