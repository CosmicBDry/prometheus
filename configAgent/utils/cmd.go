package utils

import (
	"os/exec"
)

func CmdRun(cmd string) (string, error) {
	CMD := exec.Command("bash", "-c", cmd)
	result, err := CMD.CombinedOutput()
	return string(result), err
}
