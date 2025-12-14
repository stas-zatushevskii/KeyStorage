package navigator

import tea "github.com/charmbracelet/bubbletea"

type Navigator struct {
	stack []tea.Model
}

func (n *Navigator) Push(page tea.Model) {
	n.stack = append(n.stack, page)
}

func (n *Navigator) Pop() {
	if len(n.stack) > 0 {
		n.stack = n.stack[:len(n.stack)-1]
	}
}

func (n *Navigator) Current() tea.Model {
	return n.stack[len(n.stack)-1]
}

func (n *Navigator) ReplaceCurrent(page tea.Model) {
	n.stack[len(n.stack)-1] = page
}

func New() *Navigator {
	return &Navigator{
		stack: make([]tea.Model, 0),
	}
}

type NextPage struct {
	Page tea.Model
}

func NextPageCmd(page tea.Model) tea.Cmd {
	return func() tea.Msg {
		return NextPage{Page: page}
	}
}

type PreviousPage struct {
}

func PreviousPageCmd() tea.Cmd {
	return func() tea.Msg {
		return PreviousPage{}
	}
}
