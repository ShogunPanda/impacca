/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package utils

import (
	"fmt"
	"os"
	"sync"

	"github.com/ShogunPanda/tempera"
)

var outputMutex = sync.Mutex{}

// Log shows a output message
func Log(destination *os.File, message string, args ...interface{}) {
	outputMutex.Lock()
	fmt.Fprintf(destination, message, args...)
	outputMutex.Unlock()
}

// LogWithIcon shows a message with a custom icon
func LogWithIcon(destination *os.File, icon, message string, args ...interface{}) {
	message = fmt.Sprintf("%s%s\n", SpacedEmoji(icon), message)

	Log(destination, tempera.ColorizeTemplate(message), args...)
}

// Info shows a info message
func Info(message string, args ...interface{}) {
	LogWithIcon(os.Stdout, "💬", message, args...) // Emoji code: 1F4AC
}

// Success shows a success message
func Success(message string, args ...interface{}) {
	LogWithIcon(os.Stdout, "🍻", fmt.Sprintf("{green}%s{-}", message), args...) // Emoji code: 1F4AC
}

// Warn shows a warning message
func Warn(message string, args ...interface{}) {
	LogWithIcon(os.Stdout, "⚠️", fmt.Sprintf("{bold yellow}%s{-}", message), args...) // Emoji code: 26A0+FEOF
}

// Fail shows a error message
func Fail(message string, args ...interface{}) {
	LogWithIcon(os.Stderr, "❌", fmt.Sprintf("{red}%s{-}", message), args...) // Emoji code: 274C
}

// Debug shows a debug message
func Debug(message string, args ...interface{}) {
	LogWithIcon(os.Stderr, "⚙️", fmt.Sprintf("{bold ANSI:3,0,3}%s{-}", message), args...) // Emoji code: 2699+FEOF
}

// Fatal aborts the executable with a error message
func Fatal(message string, args ...interface{}) {
	Fail(message, args...)
	os.Exit(1)
}

// Complete shows a completion message.
func Complete() {
	Success("All operations completed successfully!")
}

// NotifyStep notifies about a execution of a step
func NotifyStep(showOnly bool, color, showOnlyVerb, realVerb, message string, args ...interface{}) bool {
	verb := realVerb

	if showOnly {
		verb = showOnlyVerb
	}

	if color == "" {
		color = "{bold white}"
	}

	LogWithIcon(os.Stdout, "⚙️", fmt.Sprintf("%s%s%s{-}", color, verb, message), args...) // Emoji code: 2699+FEOF
	return !showOnly
}

// NotifyExecution notifies about a execution of a operation
func NotifyExecution(showOnly bool, showOnlyVerb, realVerb, message string, args ...interface{}) bool {
	return NotifyStep(showOnly, "{bold ANSI:3,0,3}", showOnlyVerb, realVerb, message, args...)
}

// FinishStep shows a step completion message.
func FinishStep(code int) {
	color := "green"

	if code != 0 {
		color = "red"
	}

	LogWithIcon(os.Stdout, "⚙️", fmt.Sprintf("{%s}Exited with status %d{-}", color, code)) // Emoji code: 2699+FEOF
}
