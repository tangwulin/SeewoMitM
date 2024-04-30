package helper

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func FindPidByName(name string) (int, error) {
	res := -1
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("tasklist | findstr %s", name)
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.Stdout = &outBytes
	err := cmd.Run()
	if err != nil {
		return -1, err
	}

	resStr := outBytes.String()

	if len(resStr) == 0 {
		return -1, errors.New("no such process")
	}

	parts := strings.Fields(resStr)
	if len(parts) > 1 {
		pid, err := strconv.Atoi(parts[1])
		if err != nil {
			return -1, err
		}
		res = pid
	}
	return res, nil
}

func StartProcess(path string) error {
	cmd := exec.Command(path)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
