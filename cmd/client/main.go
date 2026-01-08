package main

import (
	"client/internal/app"
	nav "client/internal/navigator"
	"client/internal/pages/auth"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

func main() {
	navigator := nav.New()

	client := resty.New()

	navigator.Push(auth.NewPage(client))

	root := app.New(navigator, &app.UserState{})

	p := tea.NewProgram(root)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
