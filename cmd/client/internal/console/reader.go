package console

import (
	"context"
	"os"
)

func ConsoleReader(ctx context.Context, ch chan []byte) error {
	b := make([]byte, 5)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err := os.Stdin.Read(b)

			data := make([]byte, 0, 5)
			for i := 0; i < 5; i++ {
				if b[i] != '\x00' {
					data = append(data, b[i])
				}
			}

			if len(data) == 0 {
				continue
			}
			if err != nil {
				return err
			}
			ch <- data
		}
	}

}
