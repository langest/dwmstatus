package main

import (
	"os/exec"
	"bytes"
	"strings"
)

func getKeyboardLayout() string {
	cmd := exec.Command("setxkbmap", "-query")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err.Error()
	}
	return strings.Fields(out.String())[5]
}
