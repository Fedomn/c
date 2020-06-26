package main

func main() {
	selectedCommandChan := make(chan Cmd)
	uiList := NewUIList(LoadCommands(), selectedCommandChan)
	uiList.registerRsyncUploader(RsyncPlugin{})

	go uiList.ListenEvents()

	for command := range selectedCommandChan {
		ExecCommand(command)
	}
}
