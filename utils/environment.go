/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type npmPackageJSON struct {
	Private bool `json:"private"`
}

var versionMatcher = regexp.MustCompile("^(v(?:-?))")

const (
	// PlainPackageManager releases using Git
	PlainPackageManager int = iota
	// NpmPackageManager releases using npm
	NpmPackageManager
	// GemPackageManager release using "rake release" task
	GemPackageManager
)

// DetectPackageManager detects which kind of release we have to use
func DetectPackageManager() int {
	cwd, _ := os.Getwd()

	if _, err := os.Stat(filepath.Join(cwd, "package.json")); !os.IsNotExist(err) {
		// Parse the package.json
		var parsed npmPackageJSON
		rawConfiguration, err := ioutil.ReadFile(filepath.Join(cwd, "package.json"))

		if err == nil {
			err = json.Unmarshal(rawConfiguration, &parsed)
		}

		// If the package.json file is marked as private, treat as PlainRelease
		if err != nil || !parsed.Private {
			return NpmPackageManager
		}
	} else if specs, err := filepath.Glob(filepath.Join(cwd, "*.gemspec")); err == nil && len(specs) > 0 {
		return GemPackageManager
	}

	return PlainPackageManager
}
