package upload

import (
	"client/internal/app"
	nav "client/internal/navigator"
	"strings"

	errorPage "client/internal/pages/error"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	app       *app.Ctx
	inputs    []textinput.Model
	focus     int
	uploading bool
}

type uploadDoneMsg struct {
	err error
}

func NewPage(app *app.Ctx) tea.Model {
	path := textinput.New()
	path.Placeholder = "/path/to/file.ext"
	path.Prompt = "File path: "
	path.CharLimit = 512
	path.Focus()

	return &Model{
		app:    app,
		inputs: []textinput.Model{path},
		focus:  0,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	submitIndex := len(m.inputs)

	switch x := msg.(type) {

	case uploadDoneMsg:
		m.uploading = false
		if x.err != nil {
			return m, nav.NextPageCmd(errorPage.New(x.err))
		}
		return m, nav.PreviousPageCmd()

	case tea.KeyMsg:
		switch x.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			if m.uploading {
				return m, nil
			}
			return m, nav.PreviousPageCmd()

		case "up", "k":
			if m.uploading {
				return m, nil
			}
			if m.focus > 0 {
				m.setFocus(m.focus - 1)
			}
			return m, nil

		case "down", "j":
			if m.uploading {
				return m, nil
			}
			if m.focus < submitIndex {
				m.setFocus(m.focus + 1)
			}
			return m, nil

		case "enter":
			if m.uploading {
				return m, nil
			}

			// Нажали на submit
			if m.focus == submitIndex {
				m.uploading = true

				filePath := strings.TrimSpace(m.inputs[0].Value())

				err := UploadFileObj(m.app, filePath)
				if err != nil {
					return m, nav.NextPageCmd(errorPage.New(err))
				}

				return m, nav.PreviousPageCmd()
			}

			if m.focus < submitIndex {
				m.setFocus(m.focus + 1)
			}
			return m, nil
		}
	}

	// Пока идёт выгрузка — поля не редактируем
	if m.uploading {
		return m, nil
	}

	// Обновляем активный input
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
	b.WriteString("Upload file\n\n")

	// поля
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// submit “кнопка”
	if m.focus == len(m.inputs) {
		b.WriteString("> [ Upload ]\n")
	} else {
		b.WriteString("  [ Upload ]\n")
	}

	b.WriteString("\n")

	if m.uploading {
		b.WriteString("Uploading... please wait\n")
		b.WriteString("(Esc/back disabled while uploading)\n")
	} else {
		b.WriteString("[↑/↓] switch   [Enter] select   [tab] back   [q] quit\n")
	}

	return b.String()
}
