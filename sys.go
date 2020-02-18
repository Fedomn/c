package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/fatih/color"
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
	fmt.Println(color.RedString("Execute %s :", cmd.Name), color.GreenString(cmd.Cmd))
}
