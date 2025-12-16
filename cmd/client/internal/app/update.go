package app

import (
	nav "client/internal/navigator"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case nav.NextPage:
		m.nav.Push(msg.Page)
		return m, nil
	case nav.PreviousPage:
		m.nav.Pop()
		return m, nil
	}
	update, cmd := m.nav.Current().Update(msg)
	m.nav.ReplaceCurrent(update)

	return m, cmd
}
