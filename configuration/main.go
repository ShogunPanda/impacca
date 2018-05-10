/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/ShogunPanda/fishamnium/console"
)

type commitMessages struct {
	Versioning string `json:"versioning"`
	Changelog  string `json:"changelog"`
}

// Configuration represents the Impacca configuration
type Configuration struct {
	CommitMessages commitMessages `json:"commitMessages"`
}

func loadConfiguration() Configuration {
	var configuration = defaultConfiguration

	// First of all, try to load the file starting from the current folder and traversing up to root - Then trying with the home folder
	visitedFolders := make(map[string]bool)
	var folders []string
	pwd, _ := os.Getwd()
	home := os.Getenv("HOME")

	for _, currentFolder := range []string{pwd, home} {
		for currentFolder != "" {
			if visitedFolders[currentFolder] {
				break
			}

			folders = append(folders, currentFolder)
			visitedFolders[currentFolder] = true

			// Go to the parent folder
			if currentFolder == "/" {
				currentFolder = ""
			} else {
				currentFolder = path.Dir(currentFolder)
			}
		}
	}

	// Check in which folder the file exists
	var configurationPath string
	for _, folder := range folders {
		tempPath := path.Join(folder, ".impacca.json")

		if _, err := os.Stat(tempPath); err == nil {
			configurationPath = tempPath
			break
		}
	}

	// Parse JSON, if any found
	if configurationPath != "" {
		rawConfiguration, err := ioutil.ReadFile(configurationPath)

		if err == nil {
			err = json.Unmarshal(rawConfiguration, &configuration)
		}

		if err != nil {
			console.Warn("The configuration file {yellow|bold}%s{-} is not a valid JSON file. Ignoring it.", configurationPath)
		}
	}

	return configuration
}

var defaultConfiguration = Configuration{
	CommitMessages: commitMessages{Versioning: "Version %s.", Changelog: "Updated CHANGELOG.md."},
}

// Current is the current Impacca configuration
var Current = loadConfiguration()
