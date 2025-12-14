package main

import (
	"client/internal/app"
	nav "client/internal/navigator"
	"client/internal/pages/mainPage"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	navigator := nav.New()

	navigator.Push(mainPage.New())

	root := app.New(navigator, &app.UserState{})

	p := tea.NewProgram(root)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
