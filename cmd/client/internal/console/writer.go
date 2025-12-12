package console

func WriteScreen(console *Console, ch chan []byte) {
	for {
		select {
		case <-ch:
			console.WriteContent()
		}
	}
}
