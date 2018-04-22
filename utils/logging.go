/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package utils

import (
	"fmt"
	"os"

	"github.com/ShogunPanda/tempera"
)

// Info shows a info message
func Info(message string, args ...interface{}) {
	message = tempera.ColorizeTemplate(fmt.Sprintf("💬%s%s\n{-}", emojiSpacer, message)) // Emoji code: 1F4AC
	fmt.Fprintf(os.Stdout, message, args...)
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

	message = tempera.ColorizeTemplate(fmt.Sprintf("⚙️%s%s%s%s\n{-}", color, emojiSpacer, verb, message)) // Emoji code: 1F4AC
	fmt.Fprintf(os.Stdout, message, args...)

	return !showOnly
}

// NotifyExecution notifies about a execution of a operation
func NotifyExecution(showOnly bool, showOnlyVerb, realVerb, message string, args ...interface{}) bool {
	return NotifyStep(showOnly, "{bold ANSI:3,0,3}", showOnlyVerb, realVerb, message, args...)
}

// Success shows a success message.
func Success(message string, args ...interface{}) {
	message = tempera.ColorizeTemplate(fmt.Sprintf("🍻%s{green}%s\n{-}", emojiSpacer, message)) // Emoji code: 1F37B
	fmt.Fprintf(os.Stdout, message, args...)
}

// Fail shows a error message.
func Fail(message string, args ...interface{}) {
	message = tempera.ColorizeTemplate(fmt.Sprintf("❌%s{red}%s\n{-}", emojiSpacer, message)) // Emoji code: 274C

	fmt.Fprintf(os.Stderr, message, args...)
}

// Fatal aborts the executable with a error message.
func Fatal(message string, args ...interface{}) {
	Fail(message, args...)
	os.Exit(1)
}

// Complete shows a completion message.
func Complete() {
	Success("All operations completed successfully!")
}

// FinishStep shows a step completion message.
func FinishStep(code int) {
	color := "green"

	if code != 0 {
		color = "red"
	}

	message := tempera.ColorizeTemplate(fmt.Sprintf("⚙️%s{%s}Exited with status %d.\n{-}", emojiSpacer, color, code)) // Emoji code: 1F4AC
	fmt.Fprintf(os.Stdout, message)
}
