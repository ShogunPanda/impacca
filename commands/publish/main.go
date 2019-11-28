/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package publish

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/ShogunPanda/impacca/utils"
	"github.com/spf13/cobra"
)

// InitCLI initializes the CLI
func InitCLI() *cobra.Command {
	cmd := &cobra.Command{
		Use: "publish <version> [changes...]", Aliases: []string{"p"}, Short: "Publishes a new version.", Args: cobra.MinimumNArgs(1), Run: publish,
	}

	cmd.Flags().BoolP("private", "p", false, "Use private scope when possible.")
	cmd.Flags().StringP("remote", "r", "origin", "The git remote name.")
	cmd.Flags().StringP("token", "t", "", "The GitHub API token.")
	cmd.Flags().BoolP("skip-changelog", "c", false, "Do not update the CHANGELOG.md file.")
	cmd.Flags().BoolP("skip-release", "R", false, "Do not update GitHub releases.")

	return cmd
}

func detectNewVersion(currentVersion *semver.Version) *semver.Version {
	changes := utils.ListChanges(currentVersion.String(), "")
	newVersion := "patch"

	if len(changes) == 0 {
		utils.Fatal("Cannot detect the new version: no changes found.")
	}

	for _, change := range changes {
		if strings.HasSuffix(change.Type, "!") || strings.Index(change.Message, "\nBREAKING CHANGE: ") != -1 {
			newVersion = "major"
			break
		} else if change.Hash == "feat" {
			newVersion = "minor"
		}
	}

	return utils.ChangeVersion(currentVersion, newVersion)
}

func publishNpmPackage(newVersion, currentVersion *semver.Version, private, dryRun bool) {
	access := "public"

	if private {
		access = "restricted"
	}

	utils.NotifyStep(dryRun, "", "Will update", "Updating", " the version to {primary}%s{-} ...", newVersion.String())
	utils.UpdateNpmVersion(newVersion, currentVersion, true, false, dryRun)

	if utils.NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}npm publish --access %s{-} ...", access) {
		result := utils.Execute(true, "npm", "publish", fmt.Sprintf("--access %s", access))
		result.Verify("npm", "Cannot publish the package")
	}
}

func publishGem(newVersion, currentVersion *semver.Version, dryRun bool) {
	utils.NotifyStep(dryRun, "", "Will update", "Updating", " the version to {primary}%s{-} ...", newVersion.String())
	utils.UpdateGemVersion(newVersion, currentVersion, true, false, dryRun)

	if utils.NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}rake release{-} ...") {
		result := utils.Execute(true, "rake", "release")
		result.Verify("rake", "Cannot publish the gem")
	}
}

func publishPlain(newVersion, currentVersion *semver.Version, dryRun bool) {
	utils.NotifyStep(dryRun, "", "Will update", "Updating", " the version to {primary}%s{-} ...", newVersion.String())
	utils.UpdateVersion(newVersion, currentVersion, dryRun)

	if utils.NotifyExecution(dryRun, "Will push", "Pushing", " commits ...") {
		result := utils.Execute(true, "git", "push")
		result.Verify("git", "Cannot push commits")
	}

	if utils.NotifyExecution(dryRun, "Will push", "Pushing", " tags ...") {
		result := utils.Execute(true, "git", "push", "--force", "--tags")
		result.Verify("git", "Cannot push tags")
	}
}

func publish(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	skipChangelog, _ := cmd.Flags().GetBool("skip-changelog")
	skipRelease, _ := cmd.Flags().GetBool("skip-release")
	private, _ := cmd.Flags().GetBool("private")
	remote, _ := cmd.Flags().GetString("remote")
	repository := utils.DetectGithubRepository(remote, true)
	token, _ := cmd.Flags().GetString("token")

	rawChanges := args[1:]
	currentVersion := utils.GetCurrentVersion()
	var newVersion *semver.Version

	if args[0] == "auto" {
		newVersion = detectNewVersion(currentVersion)
	} else {
		newVersion = utils.ChangeVersion(currentVersion, args[0])
	}

	if !dryRun {
		utils.GitMustBeClean("perform the publishing")
	}

	if !skipRelease && repository != "" && token == "" {
		utils.Fatal("In order to publish with a related GitHub release, you must provide a GitHub API token.")
	}

	if !skipChangelog && utils.NotifyStep(dryRun, "", "Will update", "Updating", " CHANGELOG.md file ...") {
		changes := make([]utils.Change, 0)

		if len(rawChanges) == 0 {
			changes = utils.ListChanges(currentVersion.String(), "")
		} else {
			for _, c := range rawChanges {
				changes = append(changes, utils.Change{Hash: "", Message: c})
			}
		}

		utils.SaveChanges(newVersion, currentVersion, changes, dryRun)
	}

	switch utils.DetectPackageManager() {
	case utils.NpmPackageManager:
		publishNpmPackage(newVersion, currentVersion, private, dryRun)
	case utils.GemPackageManager:
		publishGem(newVersion, currentVersion, dryRun)
	default:
		publishPlain(newVersion, currentVersion, dryRun)
	}

	// Now edit the Github release, if applicable
	if !skipRelease && repository != "" {
		utils.SaveRelease(newVersion, repository, remote, token, dryRun)
	}

	// TODO@PI:

	utils.Complete()
}
