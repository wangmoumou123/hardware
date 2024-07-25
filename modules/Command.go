package modules

import (
	"os/exec"
	"strings"
)

type Command struct {
}

func getStatusOutput(command string, args ...string) (int, string) {
	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()

	var statusCode int
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			statusCode = exitError.ExitCode()
		} else {
			statusCode = -1
		}
	} else {
		statusCode = 0
	}

	return statusCode, strings.TrimSpace(string(output))
}

func (c Command) RunCmD(cmd string) (int, string) {
	sta, output := getStatusOutput(cmd)
	if sta == -1 {
		return sta, cmd + ": 未找到命令"
	}
	return sta, output
}
