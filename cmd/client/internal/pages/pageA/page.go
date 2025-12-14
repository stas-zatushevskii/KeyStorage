package pageA

import (
	nav "client/internal/navigator"
	"client/internal/pages/pageC"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cursor int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func NewPage() Model {
	return Model{cursor: 0}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w":
			return m, nav.NextPageCmd(pageC.NewPage())
		case "s":
			return m, nav.PreviousPageCmd()
		}
	}
	return m, nil
}

func (m Model) View() string {
	return "PageA"
}
