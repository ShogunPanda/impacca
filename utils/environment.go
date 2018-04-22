/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package utils

import (
	"os"
	"path/filepath"
	"regexp"
)

var versionMatcher = regexp.MustCompile("^(v(?:-?))")

const (
	// PlainRelease releases using Git
	PlainRelease int = iota
	// NpmRelease releases using npm
	NpmRelease
	// GemRelease release using "rake release" task
	GemRelease
)

// DetectRelease detects which kind of release we have to use
func DetectRelease() int {
	cwd, _ := os.Getwd()

	if _, err := os.Stat(filepath.Join(cwd, "package.json")); !os.IsNotExist(err) {
		return NpmRelease
	} else if specs, err := filepath.Glob(filepath.Join(cwd, "*.gemspec")); err == nil && len(specs) > 0 {
		return GemRelease
	}

	return PlainRelease
}
