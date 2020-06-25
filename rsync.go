package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"unicode"
)

type RsyncUploader interface {
	Upload(cmd Cmd) (string, error)
}

type RsyncPlugin struct{}

const osaScript = "osascript"
const darwin = "darwin"

var chooseFileArgs = []string{
	`-e`, `tell application "iTerm2" to activate`,
	`-e`, `tell application "iTerm2" to set thefile to choose file with prompt "Choose a file to send"`,
	`-e`, `do shell script ("echo "&(quoted form of POSIX path of thefile as Unicode text)&"")`,
}

var (
	ErrRsOs         = fmt.Errorf("rsync only supported on darwin")
	ErrRsIterm2     = fmt.Errorf("rsync only supported on iTerm2")
	ErrRsNotSSHCmd  = fmt.Errorf("rsync only supported cmd pattern: ssh -i key user@ip")
	ErrRsUserCancel = fmt.Errorf("ignore: user canceled choose file")
)

func (r RsyncPlugin) resolveSSHCmd(cmdStr string) ([]string, error) {
	// pattern: ssh -i key user@ip
	cmdFields := strings.FieldsFunc(cmdStr, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	if len(cmdFields) != 4 || strings.Join(cmdFields[:2], " ") != "ssh -i" || len(strings.Split(cmdFields[3], "@")) != 2 {
		return []string{}, ErrRsNotSSHCmd
	}

	return cmdFields, nil
}

func (r RsyncPlugin) interactFile() (string, error) {
	debug("Rsync platform: %+v", runtime.GOOS)
	if runtime.GOOS != darwin {
		return "", ErrRsOs
	}

	chooseFileOutputs, err := exec.Command(osaScript, chooseFileArgs...).CombinedOutput()
	if err != nil {
		if strings.Contains(string(chooseFileOutputs), "User canceled. (-128)") {
			return "", ErrRsUserCancel
		} else {
			return "", ErrRsIterm2
		}
	}

	chooseFilePath := strings.TrimSpace(string(chooseFileOutputs))
	debug("Rsync chooseFilePath: %s", chooseFilePath)
	return chooseFilePath, nil
}

func (r RsyncPlugin) buildRsyncCmd(cmdFields []string, chooseFilePath string) (string, error) {
	// ssh -i /key
	sshCmdStr := strings.Join(cmdFields[:3], " ")

	// user@ip
	destHost := cmdFields[3]
	destUser := strings.Split(destHost, "@")[0]
	destDir := fmt.Sprintf("/home/%s", destUser)

	// user@ip:/home/user
	destStr := fmt.Sprintf("%s:%s", destHost, destDir)

	// rsync -azP -e "ssh -i key" local_file  user@ip:/home/user
	rsyncCmdStr := fmt.Sprintf(`rsync -azP -e "%s" %s %s`, sshCmdStr, chooseFilePath, destStr)
	debug("Rsync Cmd: %s", rsyncCmdStr)

	return rsyncCmdStr, nil
}

func (r RsyncPlugin) Upload(cmd Cmd) (string, error) {
	cmdFields, err := r.resolveSSHCmd(cmd.Cmd)
	if err != nil {
		return "", err
	}

	chooseFilePath, err := r.interactFile()
	if err != nil {
		return "", err
	}

	rsyncCmd, err := r.buildRsyncCmd(cmdFields, chooseFilePath)
	if err != nil {
		return "", err
	}

	return rsyncCmd, nil
}
