package auth

import (
	"client/internal/app"
	"fmt"
	"strings"

	nav "client/internal/navigator"
	"client/internal/pages/login"
	"client/internal/pages/registration"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	appCtx *app.Ctx
	items  []string
	cursor int
}

const (
	Registration = "registration"
	Login        = "login"
)

func NewPage(app *app.Ctx) tea.Model {
	return &Model{
		items: []string{
			Registration,
			Login,
		},
		cursor: 0,
		appCtx: app,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// навигация по кнопкам
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
			return m, nil

		// выбор
		case "enter":
			switch m.items[m.cursor] {
			case Registration:
				return m, nav.NextPageCmd(registration.NewPage(m.appCtx))
			case Login:
				return m, nav.NextPageCmd(login.NewPage(m.appCtx))
			}
			return m, nil

		// optional: выйти
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("Authorisation page\n\n")
	b.WriteString("Choose action:\n\n")

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
	}

	b.WriteString("\n[↑/↓] переключение   [Enter] выбор)\n")
	return b.String()
}
