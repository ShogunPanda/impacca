/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"syscall"
)

// ExecutionResult represents a command execution result
type ExecutionResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
}

// Verify checks the command executed properly and exited with exit code 0
func (e ExecutionResult) Verify(executableName, failureMessage string) {
	if e.Error != nil {
		Fatal("%s: {errorPrimary}%s{-}", failureMessage, e.Error.Error())
	} else if e.ExitCode != 0 {
		Fatal("%s: %s failed with code {errorPrimary}%d{-}.", failureMessage, executableName, e.ExitCode)
	}
}

func wrapOutput(output string) string {
	replacer, _ := regexp.Compile("(?m)(^)")
	return replacer.ReplaceAllString(output, "⛓️\x1b[4G$1")
}

func showAndBufferOutput(source io.ReadCloser, buffer *string, destination *os.File) {
	defer source.Close()

	scanner := bufio.NewScanner(source)

	for scanner.Scan() {
		line := scanner.Text()

		if destination != nil {
			fmt.Fprintln(destination, wrapOutput(line))
		}

		*buffer += line + "\n"
	}
}

// Execute executes a command.
func Execute(showOutput bool, cmd string, args ...string) (result ExecutionResult) {
	gitCmd := exec.Command(cmd, args...)

	// Pipe stdout and stderr
	var destinationOut, destinationErr *os.File

	if showOutput {
		destinationOut = os.Stdout
		destinationErr = os.Stderr
	}

	commandStdout, _ := gitCmd.StdoutPipe()
	commandStderr, _ := gitCmd.StderrPipe()
	go showAndBufferOutput(commandStdout, &result.Stdout, destinationOut)
	go showAndBufferOutput(commandStderr, &result.Stderr, destinationErr)

	// Execute the command
	result.Error = gitCmd.Run()

	// The command exited with errors, copy the exit code
	if result.Error != nil {
		if exitError, casted := result.Error.(*exec.ExitError); casted {
			result.Error = nil // Reset the error since it just a command failure
			result.ExitCode = exitError.Sys().(syscall.WaitStatus).ExitStatus()
		}
	}

	if showOutput {
		FinishStep(result.ExitCode)
	}

	return
}

// GitMustBeClean checks that the current working copy has not uncommitted changes
func GitMustBeClean(reason string) {
	// Execute the command
	result := Execute(false, "git", "status", "--short")
	result.Verify("git", "Cannot check repository status")

	if len(result.Stdout) > 0 {
		Fatal("Cannot {errorPrimary}%s{-} as the working directory is not clean. Please commit all local changes and try again.", reason)
	}
}
