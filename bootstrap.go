// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	binaryName := "mdreview"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// 1. Check if binary exists
	_, err := os.Stat(binaryName)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Binary not found. Building mdreview for the first time...\n")
		buildCmd := exec.Command("go", "build", "-o", binaryName, ".")
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to build mdreview: %v\n", err)
			os.Exit(1)
		}
	}

	// 2. Execute the binary with all passed arguments
	cmd := exec.Command("./"+binaryName, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}
