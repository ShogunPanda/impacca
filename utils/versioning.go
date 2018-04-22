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
	"sort"
	"strings"

	"github.com/Masterminds/semver"
)

func commitVersioning(version *semver.Version, commit, tag, dryRun bool) {
	versionString := version.String()

	// Commit changes
	if commit && NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}git commit -a -m \"Version %s.\"{-} ...", versionString) {
		result := Execute(true, "git", "commit", "-a", fmt.Sprintf("--message=Version %s.", versionString))
		result.Verify("git", "Cannot commit version change")
	}

	// Tag the version
	if tag && NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}git tag -f v%s{-} ...", versionString) {
		result := Execute(true, "git", "tag", "-f", "v"+versionString)
		result.Verify("git", "Cannot tag GIT version")
	}
}

// GetVersions return all current GIT versions.
func GetVersions() semver.Collection {
	result := Execute(false, "git", "tag")
	result.Verify("git", "Cannot list GIT tags")

	var versions semver.Collection
	for _, tag := range strings.Split(strings.TrimSpace(result.Stdout), "\n") {
		if !versionMatcher.MatchString(tag) {
			continue
		}

		version, err := semver.NewVersion(versionMatcher.ReplaceAllString(tag, ""))

		if err != nil {
			Fail("Cannot parse GIT tag {errorPrimary}%s{-} as a version, will skip it: {errorPrimary}%s{-}", tag, err.Error())
			continue
		}

		versions = append(versions, version)
	}

	// Sort versions
	sort.Sort(versions)

	return versions
}

// GetCurrentVersion return the current version.
func GetCurrentVersion() *semver.Version {
	// Get the current version
	versions := GetVersions()

	if len(versions) == 0 {
		version, _ := semver.NewVersion("0.0.0")
		return version
	}

	return versions[len(versions)-1]
}

// ChangeVersion changes the current version.
func ChangeVersion(version *semver.Version, change string) *semver.Version {
	newVersion := &semver.Version{}
	var err error

	switch change {
	case "patch":
		*newVersion = version.IncPatch()
	case "minor":
		*newVersion = version.IncMinor()
	case "major":
		*newVersion = version.IncMajor()
	default:
		newVersion, err = semver.NewVersion(change)

		if err != nil {
			Fatal("Cannot parse {errorPrimary}%s{-} as a version: {errorPrimary}%s{-}", change, err.Error())
		}
	}

	return newVersion
}

// UpdateVersion updates the current version.
func UpdateVersion(newVersion, currentVersion *semver.Version, dryRun bool) {
	switch DetectRelease() {
	case NpmRelease:
		UpdateNpmVersion(newVersion, currentVersion, true, true, dryRun)
	case GemRelease:
		UpdateGemVersion(newVersion, currentVersion, true, true, dryRun)
	default:
		UpdatePlainVersion(newVersion, currentVersion, true, true, dryRun)
	}
}

// UpdateNpmVersion updates the current version using NPM.
func UpdateNpmVersion(newVersion, currentVersion *semver.Version, commit, tag, dryRun bool) {
	versionString := newVersion.String()
	if !NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}npm version %s{-} ...", versionString) {
		return
	}

	result := Execute(true, "npm", "version", versionString, "--message=Version %s.")
	result.Verify("npm", "Cannot update NPM version")
}

// UpdateGemVersion updates the current version by manipulating the version file.
func UpdateGemVersion(newVersion, currentVersion *semver.Version, commit, tag, dryRun bool) {
	cwd, _ := os.Getwd()
	files, _ := filepath.Glob(filepath.Join(cwd, "*/*/version.rb"))

	if len(files) != 1 {
		Fatal("Found no or more than one possible gem version files.")
	}

	// Open the version file
	versionFile := files[0]
	rawVersionContents, err := ioutil.ReadFile(versionFile)

	if err != nil {
		Fatal("Cannot read gem version file {errorPrimary}%s{-}: {errorPrimary}%s{-}", versionFile, err.Error())
	}

	if !dryRun {
		versionContents := string(rawVersionContents)

		// Replace contents
		versionContents = regexp.MustCompile("(?m)^(?:(\\s*MAJOR)\\s*=\\s*\\d+)$").ReplaceAllString(versionContents, fmt.Sprintf("$1 = %d", newVersion.Major()))
		versionContents = regexp.MustCompile("(?m)^(?:(\\s*MINOR)\\s*=\\s*\\d+)$").ReplaceAllString(versionContents, fmt.Sprintf("$1 = %d", newVersion.Minor()))
		versionContents = regexp.MustCompile("(?m)^(?:(\\s*PATCH)\\s*=\\s*\\d+)$").ReplaceAllString(versionContents, fmt.Sprintf("$1 = %d", newVersion.Patch()))

		err := ioutil.WriteFile(versionFile, []byte(versionContents), 0644)

		if err != nil {
			Fatal("Cannot update gem version file {errorPrimary}%s{-}: {errorPrimary}%s{-}", versionFile, err.Error())
		}
	}

	commitVersioning(newVersion, commit, tag, dryRun)
}

// UpdatePlainVersion updates the current version according to a plain managament.
func UpdatePlainVersion(newVersion, currentVersion *semver.Version, commit, tag, dryRun bool) {
	versionString := newVersion.String()

	cwd, _ := os.Getwd()
	stat, err := os.Stat(filepath.Join(cwd, "Impaccafile"))

	if err == nil && stat.IsDir() == false && stat.Mode()&0111 != 0 {
		if !NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}./Impaccafile %s %s{-} ...", newVersion, currentVersion) {
			return
		}

		result := Execute(true, filepath.Join(cwd, "Impaccafile"), versionString, currentVersion.String())
		result.Verify("git", "Cannot execute the Impaccafile")

		if commit {
			if !NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}git commit -a -m \"Version %s.\"{-} ...", versionString) {
				return
			}

			result = Execute(true, "git", "commit", "-a", fmt.Sprintf("-m Version %s", versionString))
			result.Verify("Impaccafile", "Cannot commit Impaccafile changes")
		}
	}

	commitVersioning(newVersion, false, tag, dryRun)
}
