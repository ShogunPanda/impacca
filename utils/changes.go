/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/ShogunPanda/impacca/configuration"
)

// Change represents a git commit
type Change struct {
	Hash    string
	Message string
}

// ListChanges lists changes since the last version.
func ListChanges(version string) []Change {
	// Get the current version
	if version == "" {
		versions := GetVersions()
		version = versions[len(versions)-1].String()
	}

	// Get the list of changes
	result := Execute(false, "git", "log", "--format=%h %s", fmt.Sprintf("v%s..HEAD", version))
	result.Verify("git", "Cannot list GIT changes")

	rawChanges := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	changes := make([]Change, 0)

	for _, change := range rawChanges {
		if change == "" {
			continue
		}

		changeTokens := strings.SplitN(change, " ", 2)
		changes = append(changes, Change{changeTokens[0], changeTokens[1]})
	}

	return changes
}

// SaveChanges persist changes from GIT to the CHANGELOG.md file.
func SaveChanges(newVersion, currentVersion *semver.Version, changes []Change, dryRun bool) {
	cwd, _ := os.Getwd()
	changelog, err := ioutil.ReadFile(filepath.Join(cwd, "CHANGELOG.md"))

	if len(changes) == 0 {
		changes = ListChanges(currentVersion.String())
	}

	if err != nil {
		Fatal("Cannot read file {errorPrimary}CHANGELOG.md{-}: {errorPrimary}%s{-}", err.Error())
	}

	// Create the new entry
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("### %s / %s\n\n", time.Now().Format("2006-01-02"), newVersion.String()))

	for _, change := range changes {
		builder.WriteString(fmt.Sprintf("* %s\n", change.Message))
	}

	// Append the existing Changelog
	builder.WriteString("\n")
	builder.Write(changelog)

	// Save the new file
	err = ioutil.WriteFile(filepath.Join(cwd, "CHANGELOG.md"), []byte(builder.String()), 0644)

	if err != nil {
		Fatal("Cannot update file {errorPrimary}CHANGELOG.md{-}: {errorPrimary}%s{-}", err.Error())
	}

	// Commit changes
	if NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}git commit -a -m \"%s\"{-} ...", configuration.Current.CommitMessages.Changelog) {
		result := Execute(true, "git", "commit", "-a", fmt.Sprintf("--message=%s", configuration.Current.CommitMessages.Changelog))
		result.Verify("git", "Cannot commit CHANGELOG.md update")
	}
}
