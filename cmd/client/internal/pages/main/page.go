package main_page

import (
	"client/internal/app"
	"client/internal/constants"
	nav "client/internal/navigator"
	"fmt"
	"strings"

	"client/internal/pages/obj_types"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	app    *app.Ctx
	items  []string
	cursor int
}

const (
	MyStorage = "my storage"
	Upload    = "upload"
)

func NewPage(app *app.Ctx) tea.Model {
	return &Model{
		items: []string{
			MyStorage,
			Upload,
		},
		cursor: 0,
		app:    app,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

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

		case "enter":
			switch m.items[m.cursor] {

			case MyStorage:
				// LIST mode (показать списки объектов)
				return m, nav.NextPageCmd(obj_types.NewPage(m.app, constants.ModeList))

			case Upload:
				// CREATE mode (создать новый объект)
				return m, nav.NextPageCmd(obj_types.NewPage(m.app, constants.ModeCreate))
			}
			return m, nil
		case "b":
			return m, nav.PreviousPageCmd()

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("Main page\n\n")
	b.WriteString("Choose action:\n\n")

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
	}

	b.WriteString("\n[↑/↓] переключение   [Enter] выбрать   [b] назад)\n")
	return b.String()
}
