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
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/ShogunPanda/impacca/configuration"
)

var commitChecker = regexp.MustCompile("^[a-f0-9]+$")
var updateChangelogCommitFilter = regexp.MustCompile("(?i)^(?:(update(?:[ds])? changelog(?:\\.md)?(?:.)?))$")
var versionTagCommitFilter = regexp.MustCompile("(?i)^(?:version\\s+\\d+\\.\\d+\\.\\d+(?:.)?)$")

// Change represents a git commit
type Change struct {
	Hash    string
	Message string
	Type    string
}

// GetFirstCommitHash gets the first commit hash
func GetFirstCommitHash() string {
	result := Execute(false, "git", "log", "--reverse", "--format=%H")
	result.Verify("git", "Cannot get first GIT commit")

	return strings.Split(strings.TrimSpace(result.Stdout), "\n")[0]
}

// ListChanges lists changes since the last version or between specific version.
func ListChanges(version, previousVersion string) []Change {
	// Get the current version
	if version == "" {
		versions := GetVersions()
		version = versions[len(versions)-1].String()
	}

	if previousVersion == "" {
		previousVersion = "HEAD"
	}

	if version != "HEAD" && !commitChecker.MatchString(version) {
		version = fmt.Sprintf("v%s", version)
	}

	if previousVersion != "HEAD" && !commitChecker.MatchString(previousVersion) {
		previousVersion = fmt.Sprintf("v%s", previousVersion)
	}

	// Get the list of changes
	executionArgs := []string{"log", "--format=%h %s"}

	if version != "0.0.0" {
		executionArgs = append(executionArgs, fmt.Sprintf("%s...%s", previousVersion, version))
	}

	result := Execute(false, "git", executionArgs...)
	result.Verify("git", "Cannot list GIT changes")

	rawChanges := strings.Split(strings.TrimSpace(result.Stdout), "\n")
	changes := make([]Change, 0)

	for _, change := range rawChanges {
		if change == "" {
			continue
		}

		changeTokens := strings.SplitN(change, " ", 2)
		messageComponents := []string{"feat", changeTokens[1]}

		if strings.Index(messageComponents[1], ":") != -1 {
			messageComponents = strings.SplitN(changeTokens[1], ":", 2)
		} else if strings.Index(messageComponents[1], "fix") != -1 {
			messageComponents[0] = "fix"
		}

		changes = append(
			changes,
			Change{
				strings.ToLower(strings.TrimSpace(changeTokens[0])),
				strings.TrimSpace(messageComponents[1]),
				strings.ToLower(strings.TrimSpace(messageComponents[0])),
			},
		)
	}

	return changes
}

// FormatChanges formats changes to the CHANGELOG.md file format.
func FormatChanges(previous string, version *semver.Version, changes []Change, date time.Time) string {
	// Create the new entry
	var builder strings.Builder

	if !date.IsZero() {
		builder.WriteString(fmt.Sprintf("### %s / %s\n\n", date.Format("2006-01-02"), version.String()))
	}

	for _, change := range changes {
		// Filter some commits
		if updateChangelogCommitFilter.MatchString(change.Message) || versionTagCommitFilter.MatchString(change.Message) {
			continue
		}

		builder.WriteString(fmt.Sprintf("- %s: %s\n", change.Type, change.Message))
	}

	// Append the existing Changelog
	builder.WriteString("\n")
	builder.Write([]byte(previous))

	return builder.String()
}

// FormatReleaseChanges formats changes for a GitHub release.
func FormatReleaseChanges(repository string, changes []Change) string {
	// Create the new entry
	var builder strings.Builder

	for _, change := range changes {
		// Filter some commits
		if updateChangelogCommitFilter.MatchString(change.Message) || versionTagCommitFilter.MatchString(change.Message) {
			continue
		}

		builder.WriteString(
			fmt.Sprintf(
				"- %s: %s ([%s](https://github.com/%s/commit/%s))\n", 
				change.Type, change.Message, change.Hash, repository, change.Hash,
			),
		)
	}

	return builder.String()
}


// SaveChanges persist changes from GIT to the CHANGELOG.md file.
func SaveChanges(newVersion, currentVersion *semver.Version, changes []Change, dryRun bool) {
	cwd, _ := os.Getwd()
	changelog := ""
	var err error

	if _, err := os.Stat(filepath.Join(cwd, "CHANGELOG.md")); !os.IsNotExist(err) {
		rawChangelog, err := ioutil.ReadFile(filepath.Join(cwd, "CHANGELOG.md"))

		if err != nil {
			Fatal("Cannot read file {errorPrimary}CHANGELOG.md{-}: {errorPrimary}%s{-}", err.Error())
		}

		changelog = string(rawChangelog)
	}

	if len(changes) == 0 {
		changes = ListChanges(currentVersion.String(), "")
	}

	if NotifyExecution(dryRun, "Will append", "Appending", " {primary}%d{-} entries to the CHANGELOG.md file ...", len(changes)) {
		newChangelog := FormatChanges(changelog, newVersion, changes, time.Now())

		// Save the new file
		err = ioutil.WriteFile(filepath.Join(cwd, "CHANGELOG.md"), []byte(newChangelog), 0644)

		if err != nil {
			Fatal("Cannot update file {errorPrimary}CHANGELOG.md{-}: {errorPrimary}%s{-}", err.Error())
		}
	}

	// Commit changes
	message := strings.TrimSpace(configuration.Current.CommitMessages.Changelog)
	if NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}git commit --all --message \"%s\"{-} ...", message) {
		result := Execute(true, "git", "add", "CHANGELOG.md")
		result.Verify("git", "Cannot add CHANGELOG.md update to git stage area")

		result = Execute(true, "git", "commit", "--all", fmt.Sprintf("--message=%s", message))
		result.Verify("git", "Cannot commit CHANGELOG.md update")
	}
}
