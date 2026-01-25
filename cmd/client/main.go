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

	appCtx := app.NewCtx(client)

	navigator.Push(auth.NewPage(appCtx))

	root := app.New(navigator)

	p := tea.NewProgram(root)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
