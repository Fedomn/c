package main

func main() {
	selectedCommandChan := make(chan Cmd)

	go NewUIList(LoadCommands(), selectedCommandChan).ListenEvents()

	for {
		select {
		case command := <-selectedCommandChan:
			ExecCommand(command)
		}
	}
}
