package console

import (
	"bytes"
	"client/internal/console/codes/console"
	"errors"
	"fmt"
	"strings"
)

type Console struct {
	currentLineIndex int
	Lines            []*Line
}

type Line struct {
	cursorUp             string
	cursorDown           string
	textColor            string
	backgroundColor      string
	focusTextColor       string
	focusBackgroundColor string
	Position             int
	operatingMode        int // 0 one click only, 1 input Line
	buffer               bytes.Buffer
}

func New() *Console {
	return &Console{
		currentLineIndex: 0,
		Lines:            []*Line{},
	}
}

func (l *Line) coloredString() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s", l.backgroundColor, l.textColor, l.cursorUp, console.ClearLine, l.buffer.String(), l.cursorDown, console.Reset)
}

func (l *Line) focusColoredString() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s", l.focusBackgroundColor, l.focusTextColor, l.cursorUp, console.ClearLine, l.buffer.String(), l.cursorDown, console.Reset)
}

func (l *Line) string() string {
	return fmt.Sprintf("%s%s%s%s%s", l.cursorUp, console.ClearLine, l.buffer.String(), l.cursorDown, console.Reset)
}

func (c *Console) AddNewLine() {
	var line = &Line{}
	n := len(c.Lines)

	line.cursorUp = strings.Repeat(console.CursorUp, n)
	line.cursorDown = strings.Repeat(console.CursorDown, n)

	line.Position = len(c.Lines)

	c.Lines = append(c.Lines, line)
}

func (c *Console) WriteLine(linePosition int, content string) error {
	if linePosition < 0 || linePosition >= len(c.Lines) {
		return errors.New("line out of range")
	}
	c.Lines[linePosition].buffer.WriteString(content)
	return nil
}

func (c *Console) SetLineTextColor(linePosition int, color string) error {
	if linePosition < 0 || linePosition >= len(c.Lines) {
		return errors.New("line out of range")
	}
	c.Lines[linePosition].textColor = color
	return nil
}

func (c *Console) SetFocusLineTextColor(linePosition int, color string) error {
	if linePosition < 0 || linePosition >= len(c.Lines) {
		return errors.New("line out of range")
	}
	c.Lines[linePosition].focusTextColor = color
	return nil
}

func (c *Console) SetFocusLineBackgroundColor(linePosition int, color string) error {
	if linePosition < 0 || linePosition >= len(c.Lines) {
		return errors.New("line out of range")
	}
	c.Lines[linePosition].focusBackgroundColor = color
	return nil
}

func (c *Console) SetLineBackgroundColor(linePosition int, color string) error {
	if linePosition < 0 || linePosition >= len(c.Lines) {
		return errors.New("line out of range")
	}
	c.Lines[linePosition].backgroundColor = color
	return nil
}

func (c *Console) ClearConsole() {
	for i := 0; i < len(c.Lines); i++ {
		c.Lines[i].buffer.Reset()
	}
}

func (c *Console) WriteContent() {
	for i := 0; i < len(c.Lines); i++ {
		if c.Lines[i].Position == c.currentLineIndex {
			fmt.Printf(c.Lines[i].focusColoredString())
		}
		fmt.Printf(c.Lines[i].coloredString())
	}
}
