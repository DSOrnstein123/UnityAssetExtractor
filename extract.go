package main

import (
	"os"
	"os/exec"
)

func extract(filepaths []string) error {
	pythonPath := "./venv/Scripts/python.exe"
	cmd := exec.Command(pythonPath, append([]string{"extract.py"}, filepaths...)...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
