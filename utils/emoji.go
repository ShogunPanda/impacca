/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package utils

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func setTerminalMode(mode string) {
	cmd := exec.Command("/bin/stty", mode)
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
	cmd.Wait()
}

// GetEmojiWidth Detects handling of emoji
func GetEmojiWidth() int64 {
	setTerminalMode("raw")

	os.Stdout.Write([]byte("💬\x1b[6n"))
	reader := bufio.NewReader(os.Stdin)
	position, _ := reader.ReadSlice('R')

	// Set the terminal back from raw mode to 'cooked'
	setTerminalMode("-raw")

	// Delete the current line
	os.Stdout.Write([]byte("\x1b[0E\x1b[0K"))

	// Parse the position
	coordinates := strings.Split(string(position[2:len(position)-1]), ";")
	width, _ := strconv.ParseInt(coordinates[1], 0, 4)

	return width
}

// SpacedEmoji returns an emoji with a trailing space
func SpacedEmoji(emoji string) string {
	return emoji + "\x1b[4G"
}
