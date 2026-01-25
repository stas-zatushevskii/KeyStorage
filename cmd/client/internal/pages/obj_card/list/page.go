package get

import (
	"client/internal/app"
	nav "client/internal/navigator"
	"context"
	"fmt"
	"strings"
	"time"

	errorPage "client/internal/pages/error"

	get_obj "client/internal/pages/obj_card/get"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	app     *app.Ctx
	loading bool
	items   []Card
	cursor  int
}

func NewPage(app *app.Ctx) tea.Model {
	return &Model{
		app:     app,
		loading: true,
		cursor:  0,
	}
}

type listLoadedMsg struct {
	items []Card
	err   error
}

func (m Model) Init() tea.Cmd {
	// Если данные уже есть — ничего не грузим
	if !m.loading {
		return nil
	}
	return fetchListCmd(m.app)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {
	case listLoadedMsg:
		m.loading = false
		if x.err != nil {
			return m, nav.NextPageCmd(errorPage.New(x.err))
		}
		m.items = x.items
		if m.cursor >= len(m.items) {
			m.cursor = 0
		}
		return m, nil

	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "r":
			m.loading = true
			return m, fetchListCmd(m.app)

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
			if len(m.items) == 0 {
				return m, nil
			}
			selected := m.items[m.cursor]
			return m, nav.NextPageCmd(get_obj.NewPage(m.app, selected.CardID))
		case "b":
			return m, nav.PreviousPageCmd()
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("Banks\n\n")

	if m.loading {
		b.WriteString("Loading...\n\n")
		b.WriteString("[r] refresh   [q] quit\n")
		return b.String()
	}

	if len(m.items) == 0 {
		b.WriteString("(empty)\n\n")
		b.WriteString("[r] refresh   [q] quit\n")
		return b.String()
	}

	for i, it := range m.items {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		b.WriteString(fmt.Sprintf("%s%s\n", prefix, it.BankName))
	}

	b.WriteString("\n[↑/↓] переключение   [enter] открыть   [r]   обновить   [b] назад\n")
	return b.String()
}

func fetchListCmd(app *app.Ctx) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		items, err := GetCardList(ctx, app)
		return listLoadedMsg{items: items, err: err}
	}
}
