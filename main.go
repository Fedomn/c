package main

func main() {
	selectedCommandChan := make(chan Cmd)

	go NewUIList(LoadCommands(), selectedCommandChan).ListenEvents()

	for command := range selectedCommandChan {
		ExecCommand(command)
	}
}
