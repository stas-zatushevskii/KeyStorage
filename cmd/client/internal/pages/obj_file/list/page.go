package get

import (
	"client/internal/app"
	nav "client/internal/navigator"
	"context"
	"fmt"
	"strings"
	"time"

	errorPage "client/internal/pages/error"

	load_file "client/internal/pages/obj_file/load"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	app     *app.Ctx
	loading bool
	items   []File
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
	items []File
	err   error
}

func (m Model) Init() tea.Cmd {
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
			return m, nav.NextPageCmd(load_file.NewPage(m.app, selected.ID))

		case "tab":
			return m, nav.PreviousPageCmd()
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("Files\n\n")

	if m.loading {
		b.WriteString("Loading...\n\n")
		b.WriteString("[r] refresh   [b] back   [q] quit\n")
		return b.String()
	}

	if len(m.items) == 0 {
		b.WriteString("(empty)\n\n")
		b.WriteString("[r] refresh   [b] back   [q] quit\n")
		return b.String()
	}

	for i, it := range m.items {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		b.WriteString(fmt.Sprintf("%sFile name=%s\n", prefix, it.Title))
	}

	b.WriteString("\n[↑/↓] move   [enter] open   [tab] back   [q] quit\n")
	return b.String()
}

func fetchListCmd(app *app.Ctx) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		items, err := GetFileList(ctx, app)
		return listLoadedMsg{items: items, err: err}
	}
}
