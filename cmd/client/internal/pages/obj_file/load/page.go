package load

import (
	"client/internal/app"
	nav "client/internal/navigator"
	"strings"

	errorPage "client/internal/pages/error"

	tea "github.com/charmbracelet/bubbletea"
)

type downloadDoneMsg struct {
	path string
	err  error
}

type Model struct {
	app     *app.Ctx
	id      int64
	cursor  int
	loading bool
}

func NewPage(app *app.Ctx, id int64) tea.Model {
	return &Model{
		app:    app,
		id:     id,
		cursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {

	case downloadDoneMsg:
		m.loading = false
		if x.err != nil {
			return m, nav.NextPageCmd(errorPage.New(x.err))
		}
		return m, nav.PreviousPageCmd()

	case tea.KeyMsg:
		switch x.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if !m.loading && m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case "down", "j":
			if !m.loading && m.cursor < 1 {
				m.cursor++
			}
			return m, nil

		case "enter":
			if m.loading {
				return m, nil
			}

			switch m.cursor {
			case 0:
				m.loading = true
				_, err := DownloadFileByID(m.app, m.id)
				if err != nil {
					return m, nav.NextPageCmd(errorPage.New(err))
				}
				return m, nav.DoubleBackPageCmd()

			case 1:
				return m, nav.DoubleBackPageCmd()
			}

		case "tab":
			if m.loading {
				return m, nil
			}
			return m, nav.PreviousPageCmd()
		}
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString("Download file\n\n")

	download := "  [ Download ]"
	back := "  [ Back ]"

	if m.cursor == 0 {
		download = "> [ Download ]"
	}
	if m.cursor == 1 {
		back = "> [ Back ]"
	}

	b.WriteString(download + "\n")
	b.WriteString(back + "\n\n")

	if m.loading {
		b.WriteString("Downloading... please wait\n")
		b.WriteString("(navigation disabled)\n")
	} else {
		b.WriteString("[↑/↓] switch   [enter] select   [tab] back   [q] quit\n")
	}

	return b.String()
}
