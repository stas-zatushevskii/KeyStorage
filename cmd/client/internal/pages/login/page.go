package login

import (
	nav "client/internal/navigator"
	"strings"

	errorPage "client/internal/pages/error"
	mainPage "client/internal/pages/main"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type Model struct {
	client *resty.Client
	inputs []textinput.Model
	focus  int
}

func NewPage(c *resty.Client) tea.Model {
	username := textinput.New()
	username.Placeholder = "username"
	username.Prompt = "Username: "
	username.CharLimit = 64
	username.Focus()

	password := textinput.New()
	password.Placeholder = "password"
	password.Prompt = "Password: "
	password.CharLimit = 128
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '*'

	return &Model{
		inputs: []textinput.Model{username, password},
		focus:  0,
		client: c,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	submitIndex := len(m.inputs)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.focus > 0 {
				m.setFocus(m.focus - 1)
			}
			return m, nil

		case "down", "j":
			if m.focus < submitIndex {
				m.setFocus(m.focus + 1)
			}
			return m, nil

		case "enter":
			if m.focus == submitIndex {
				username := strings.TrimSpace(m.inputs[0].Value())
				password := m.inputs[1].Value()
				tokens, err := Login(m.client, username, password)
				if err != nil {
					return m, nav.NextPageCmd(errorPage.New(err))
				}

				return m, nav.NextPageCmd(mainPage.NewPage(m.client, tokens))
			}

			if m.focus < submitIndex {
				m.setFocus(m.focus + 1)
			}
			return m, nil

		case "esc":
			m.blurAll()
			return m, nil

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.focus < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) setFocus(idx int) {
	m.blurAll()
	m.focus = idx
	if m.focus < len(m.inputs) {
		m.inputs[m.focus].Focus()
	}
}

func (m *Model) blurAll() {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("Login form\n\n")

	// поля
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// submit “кнопка”
	if m.focus == len(m.inputs) {
		b.WriteString("> [ Submit ]\n")
	} else {
		b.WriteString("  [ Submit ]\n")
	}

	b.WriteString("\n(↑/↓ переключение, Enter выбрать)\n")
	return b.String()
}
