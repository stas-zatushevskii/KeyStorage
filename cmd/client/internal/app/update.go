package app

import (
	nav "client/internal/navigator"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case nav.NextPage:
		m.nav.Push(msg.Page)
		return m, m.nav.Current().Init()
	case nav.PreviousPage:
		m.nav.Pop()
		return m, m.nav.Current().Init()
	case nav.DoubleBackPage:
		m.nav.Pop()
		m.nav.Pop()
		return m, m.nav.Current().Init()
	}
	update, cmd := m.nav.Current().Update(msg)
	m.nav.ReplaceCurrent(update)

	return m, cmd
}
