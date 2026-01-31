package get

import (
	"client/internal/app"
	nav "client/internal/navigator"
	"context"
	"fmt"
	"time"

	errorPage "client/internal/pages/error"

	tea "github.com/charmbracelet/bubbletea"
)

type textLoadedMsg struct {
	item *Account
	err  error
}

type Model struct {
	app     *app.Ctx
	id      int64
	loading bool
	item    *Account
	err     error
}

func NewPage(app *app.Ctx, id int64) tea.Model {
	return &Model{
		app:     app,
		id:      id,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchTextCmd(m.app, m.id)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch x := msg.(type) {

	case textLoadedMsg:
		m.loading = false
		if x.err != nil {
			return m, nav.NextPageCmd(errorPage.New(x.err))
		}
		m.item = x.item
		return m, nil

	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, nav.PreviousPageCmd()
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.loading {
		return "Account\n\nLoading...\n"
	}

	if m.item == nil {
		return "Account\n\nNot found\n"
	}

	return fmt.Sprintf(
		"Account\n\n"+
			"Service name: %s\n\n"+
			"Username: %s\n\n"+
			"Password: %s\n\n"+
			"esc назад\n",
		m.item.ServiceName,
		m.item.Username,
		m.item.Password,
	)
}

func fetchTextCmd(app *app.Ctx, id int64) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		item, err := GetAccountByID(ctx, app, id)
		return textLoadedMsg{
			item: item,
			err:  err,
		}
	}
}
