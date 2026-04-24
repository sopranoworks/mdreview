// +build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	binaryName := "mdreview"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// 1. Check if binary exists and is up to date
	binaryStat, err := os.Stat(binaryName)
	needsBuild := os.IsNotExist(err)

	if !needsBuild {
		// Check if any .go files are newer than the binary
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".go" && path != "bootstrap.go" {
				if info.ModTime().After(binaryStat.ModTime()) {
					needsBuild = true
					fmt.Fprintf(os.Stderr, "Source file %s is newer than binary. Rebuilding...\n", path)
					return fmt.Errorf("rebuild needed")
				}
			}
			return nil
		})
	}

	if needsBuild {
		fmt.Fprintf(os.Stderr, "Building mdreview...\n")
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
