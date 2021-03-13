package gradelib

import (
	"os"
	"os/exec"
)

func RunCmd(executable string, arguments ...string) string {
	cmd := exec.Command(executable, arguments...)
	cmd.Dir, _ = os.Getwd()
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}
