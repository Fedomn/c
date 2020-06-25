package main

import (
	"fmt"
	"strings"
	"testing"
)

const chooseFilePath = "/fake/path"

var rs = RsyncPlugin{}

func TestDetectInvalidCmd(t *testing.T) {
	var tests = []struct {
		cmdStr string
		wat    error
	}{
		{"", ErrRsNotSSHCmd},
		{"xxx", ErrRsNotSSHCmd},
		{"ssh user@ip", ErrRsNotSSHCmd},
		{"ssh  -i key  user-ip", ErrRsNotSSHCmd},
		{"ssh  -i key  user@ip", nil},
	}
	for _, tt := range tests {
		_, got := rs.resolveSSHCmd(tt.cmdStr)
		msg := fmt.Sprintf("cmdStr: %s", tt.cmdStr)
		Equals(t, msg, tt.wat, got)
	}
}

func TestBuildRsyncCmd(t *testing.T) {
	var tests = []struct {
		cmdStr []string
		wat    string
	}{
		{strings.Split("ssh -i key user@ip", " "), `rsync -azP -e "ssh -i key" /fake/path user@ip:/home/user`},
	}
	for _, tt := range tests {
		got, _ := rs.buildRsyncCmd(tt.cmdStr, chooseFilePath)
		msg := fmt.Sprintf("cmdStr: %s", tt.cmdStr)
		Equals(t, msg, tt.wat, got)
	}
}
