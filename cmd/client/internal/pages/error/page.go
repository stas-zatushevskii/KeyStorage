package auth_page

import (
	"strings"

	nav "client/internal/navigator"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	err error
}

func New(err error) tea.Model {
	return &Model{
		err: err,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			return m, nav.PreviousPageCmd()

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("Something went wrong\n\n")
	b.WriteString(m.err.Error())
	b.WriteString("\n(OK Enter)\n")
	return b.String()
}
