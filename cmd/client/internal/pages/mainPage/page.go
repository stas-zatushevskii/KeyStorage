package mainPage

import (
	nav "client/internal/navigator"
	"client/internal/pages/pageA"
	"client/internal/pages/pageB"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	items  []string
	cursor int
}

const (
	ActionPageA = "pageA"
	ActionPageB = "pageB"
)

func New() tea.Model {
	return &Model{
		items: []string{
			ActionPageA,
			ActionPageB,
		},
		cursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w":
			return m, nav.NextPageCmd(pageA.NewPage())
		case "s":
			return m, nav.NextPageCmd(pageB.NewPage())
		}
	}
	return m, nil
}

func (m Model) View() string {
	return "Page main"
}
