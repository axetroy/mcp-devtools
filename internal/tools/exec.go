package tools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ExecuteCommand executes a shell command
// Note: This uses simple command splitting and does not support quoted arguments with spaces.
// For complex shell syntax, the shell should be explicitly invoked (e.g., "sh -c <command>").
func ExecuteCommand(args map[string]interface{}) (interface{}, error) {
	command, ok := args["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command parameter is required")
	}

	workDir, _ := args["workdir"].(string)

	// Use shell for command execution to properly handle quotes and complex syntax
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	if workDir != "" {
		cmd.Dir = workDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := fmt.Sprintf("STDOUT:\n%s\n\nSTDERR:\n%s\n", stdout.String(), stderr.String())

	if err != nil {
		result += fmt.Sprintf("\nError: %v", err)
	}

	return result, nil
}

// GetEnvironment gets environment variables
func GetEnvironment(args map[string]interface{}) (interface{}, error) {
	envVars := os.Environ()
	return strings.Join(envVars, "\n"), nil
}

// GetWorkingDirectory gets the current working directory
func GetWorkingDirectory(args map[string]interface{}) (interface{}, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	return dir, nil
}
