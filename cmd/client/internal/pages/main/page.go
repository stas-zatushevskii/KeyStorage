package auth

import (
	domain "client/internal/domain/token"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type Model struct {
	client *resty.Client
	token  *domain.Token
	items  []string
	cursor int
}

const (
	MyStorage = "my storage"
	Upload    = "upload new file"
)

func NewPage(c *resty.Client, tokens *domain.Token) tea.Model {
	return &Model{
		client: c,
		token:  tokens,
		items: []string{
			MyStorage,
			Upload,
		},
		cursor: 0,
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
				return m, nil // todo: page MyStorage
			case Upload:
				return m, nil // todo: page Upload
			}
			return m, nil

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
	// fixme
	b.WriteString(fmt.Sprintf("YOUR JWT: %s\n\n", m.token.GetJWTToken()))

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, item))
	}

	b.WriteString("\n(↑/↓ + Enter)\n")
	return b.String()
}
