package console

import (
	"client/internal/console/codes/keys"
)

func (c *Console) InputHandler(b []byte) int {
	if c.Lines[c.currentLineIndex].operatingMode == 0 {
		return singleInputHandler(b, c)
	}
	return multiInputHandler(b, c)

}

func singleInputHandler(b []byte, console *Console) int {
	if string(b) == keys.InputUp {
		console.currentLineIndex = (console.currentLineIndex + 1) % len(console.Lines)

		return console.currentLineIndex
	}
	if string(b) == keys.InputDown {
		if console.currentLineIndex > 0 {
			console.currentLineIndex = (console.currentLineIndex - 1) % len(console.Lines)
			return console.currentLineIndex
		} else {
			return console.currentLineIndex
		}
	}
	if string(b) == keys.InputInsert {
		return console.currentLineIndex
	}

	return console.currentLineIndex
}

func multiInputHandler(b []byte, console *Console) int {
	if string(b) == keys.InputUp {
		console.currentLineIndex++
	}
	if string(b) == keys.InputDown {
		if console.currentLineIndex > 0 {
			console.currentLineIndex--
		}
	}
	if string(b) == keys.InputInsert {
		console.Lines[console.currentLineIndex].buffer.Write(b)
	}

	if string(b) == keys.InputDelete {
		if console.Lines[console.currentLineIndex].buffer.Len() > 0 {
			console.Lines[console.currentLineIndex].buffer.Truncate(len(b) - 1)
		}
	}
	return console.currentLineIndex

}
