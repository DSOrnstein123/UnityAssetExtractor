package main

import (
	"os"
	"os/exec"
)

func decrypt(filepaths []string) error {
	cmd := exec.Command(config.DecryptScript, filepaths...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
