package app

func (m *Model) View() string {
	return m.nav.Current().View()
}
