package main

import (
	"client/internal/console"
	colors "client/internal/console/codes/colors"
	"context"
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	hideCursor = "\x1b[?25l"
	showCursor = "\x1b[?25h"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print(hideCursor)

	c := console.New()

	c.AddNewLine() // 0
	c.AddNewLine() // 1
	c.AddNewLine() // 2
	c.AddNewLine() // 3

	c.WriteLine(0, "hello from line one")
	c.WriteLine(1, "hello from line two")
	c.WriteLine(2, "hello from line three")

	c.SetLineTextColor(0, colors.Green)
	c.SetLineTextColor(1, colors.Red)
	c.SetLineTextColor(2, colors.Cyan)

	c.SetFocusLineTextColor(0, colors.Red)
	c.SetFocusLineTextColor(1, colors.Green)
	c.SetFocusLineTextColor(2, colors.Black)

	c.SetLineBackgroundColor(0, colors.BgBlack)
	c.SetLineBackgroundColor(1, colors.BgWhite)
	c.SetLineBackgroundColor(2, colors.BgGreen)

	c.SetFocusLineBackgroundColor(0, colors.White)
	c.SetFocusLineBackgroundColor(1, colors.BgBlack)
	c.SetFocusLineBackgroundColor(2, colors.BgYellow)

	ch := make(chan []byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go console.ConsoleReader(ctx, ch)
	go console.WriteScreen(c, ch)
	select {}
}
