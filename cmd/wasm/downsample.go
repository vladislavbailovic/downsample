package main

func main() {
	done := make(chan bool)

	initGui()

	for {
		select {
		case <-done:
			return
		}
	}
}
