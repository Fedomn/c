package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func ExecCommand(cmd Cmd) {
	if cmd.Cmd == "" {
		os.Exit(0)
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		panic(err)
	}
	args := []string{"bash", "-c", cmd.Cmd}
	env := os.Environ()

	printCmdInfo(cmd)

	execErr := syscall.Exec(bash, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func printCmdInfo(cmd Cmd) {
	fmt.Printf("\033[0;31mExecute %s : \033[0;39m", cmd.Name)
	fmt.Printf("\033[0;32m%s\033[0;39m\n", cmd.Cmd)
}
