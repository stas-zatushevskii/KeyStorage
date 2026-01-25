package obj_types

import (
	"client/internal/constants"
	"fmt"
	"strings"

	"client/internal/app"
	nav "client/internal/navigator"

	card_create "client/internal/pages/obj_card/create"
	card_list "client/internal/pages/obj_card/list"
	text_create "client/internal/pages/obj_text/create"
	text_list "client/internal/pages/obj_text/list"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	app    *app.Ctx
	mode   constants.Mode
	cursor int
	items  []constants.ObjType
}

func NewPage(app *app.Ctx, mode constants.Mode) tea.Model {
	return &Model{
		app:    app,
		mode:   mode,
		cursor: 0,
		items:  []constants.ObjType{constants.Text, constants.Account, constants.File, constants.Bank},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {

	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

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

		case "b":
			return m, nav.PreviousPageCmd()

		case "enter":
			if len(m.items) == 0 {
				return m, nil
			}

			obj := m.items[m.cursor]

			// В зависимости от режима — уходим либо на create, либо на list
			switch m.mode {
			case constants.ModeCreate:
				switch obj {
				case constants.Text:
					// return m, nav.NextPageCmd(upload_text_obj.NewPage(m.app))
					return m, nav.NextPageCmd(text_create.NewPage(m.app))
				case constants.Account:
					// return m, nav.NextPageCmd(create_account.NewPage(m.app))
					return m, nav.NextPageCmd(todoPage("TODO: go to CREATE Account"))
				case constants.File:
					// return m, nav.NextPageCmd(upload_file_obj.NewPage(m.app))
					return m, nav.NextPageCmd(todoPage("TODO: go to CREATE File"))
				case constants.Bank:
					// return m, nav.NextPageCmd(create_bank.NewPage(m.app))
					return m, nav.NextPageCmd(card_create.NewPage(m.app))
				default:
					panic("unhandled default case")
				}

			case constants.ModeList:
				switch obj {
				case constants.Text:
					// return m, nav.NextPageCmd(text_list.NewPage(m.app))
					return m, nav.NextPageCmd(text_list.NewPage(m.app))
				case constants.Account:
					// return m, nav.NextPageCmd(account_list.NewPage(m.app))
					return m, nav.NextPageCmd(todoPage("TODO: go to LIST Account"))
				case constants.File:
					// return m, nav.NextPageCmd(file_list.NewPage(m.app))
					return m, nav.NextPageCmd(todoPage("TODO: go to LIST File"))
				case constants.Bank:
					// return m, nav.NextPageCmd(bank_list.NewPage(m.app))
					return m, nav.NextPageCmd(card_list.NewPage(m.app))
				default:
					panic("unhandled default case")
				}
			default:
				panic("unhandled default case")
			}

			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString("Choose object type\n")
	b.WriteString(fmt.Sprintf("Mode: %s\n\n", m.mode.String()))

	for i, it := range m.items {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		b.WriteString(prefix + it.String() + "\n")
	}

	b.WriteString("\n(↑/↓ switch, Enter select, b back)\n")
	return b.String()
}

// Заглушка-страница, чтобы код компилился, пока ты не подставишь реальные pages.
// УДАЛИШЬ потом.
type simplePage string

func todoPage(text string) tea.Model { return simplePage(text) }

func (p simplePage) Init() tea.Cmd { return nil }

func (p simplePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p simplePage) View() string {
	return string(p) + "\n\n[q] quit\n"
}
