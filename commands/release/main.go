/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package release

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/ShogunPanda/impacca/utils"
	"github.com/spf13/cobra"
)

// InitCLI initializes the CLI
func InitCLI() *cobra.Command {
	cmd := &cobra.Command{
		Use: "release <version> [changes...]", Aliases: []string{"r"}, Short: "Releases a new version.", Args: cobra.MinimumNArgs(1), Run: release,
	}

	cmd.Flags().BoolP("skip-changelog", "c", false, "Do not update the CHANGELOG.md file.")
	cmd.Flags().BoolP("private", "p", false, "Use private scope when possible.")

	return cmd
}

func releaseNpmPackage(newVersion, currentVersion *semver.Version, private, dryRun bool) {
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

func releaseGem(newVersion, currentVersion *semver.Version, dryRun bool) {
	utils.NotifyStep(dryRun, "", "Will update", "Updating", " the version to {primary}%s{-} ...", newVersion.String())
	utils.UpdateGemVersion(newVersion, currentVersion, true, false, dryRun)

	if utils.NotifyExecution(dryRun, "Will execute", "Executing", ": {primary}rake release{-} ...") {
		result := utils.Execute(true, "rake", "release")
		result.Verify("rake", "Cannot publish the gem")
	}
}

func releasePlain(newVersion, currentVersion *semver.Version, dryRun bool) {
	utils.NotifyStep(dryRun, "", "Will update", "Updating", " the version to {primary}%s{-} ...", newVersion.String())
	utils.UpdateVersion(newVersion, currentVersion, dryRun)

	if utils.NotifyExecution(dryRun, "Will push", "Pushing", " commits ...") {
		result := utils.Execute(true, "git", "push")
		result.Verify("git", "Cannot push commits")
	}

	if utils.NotifyExecution(dryRun, "Will push", "Pushing", " tags ...") {
		result := utils.Execute(true, "git", "push", "-f", "--tags")
		result.Verify("git", "Cannot push tags")
	}
}

func release(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	skipChangelog, _ := cmd.Flags().GetBool("skip-changelog")
	private, _ := cmd.Flags().GetBool("private")

	rawChanges := args[1:]
	currentVersion := utils.GetCurrentVersion()
	newVersion := utils.ChangeVersion(currentVersion, args[0])

	changes := make([]utils.Change, 0)

	if len(rawChanges) == 0 {
		changes = utils.ListChanges(currentVersion.String())
	} else {
		for _, c := range rawChanges {
			changes = append(changes, utils.Change{Hash: "", Message: c})
		}
	}

	if !dryRun {
		utils.GitMustBeClean("perform the releasing")
	}

	if !skipChangelog && utils.NotifyStep(dryRun, "", "Will update", "Updating", " CHANGELOG.md file ...") {
		utils.SaveChanges(newVersion, currentVersion, changes, dryRun)
	}

	switch utils.DetectRelease() {
	case utils.NpmRelease:
		releaseNpmPackage(newVersion, currentVersion, private, dryRun)
	case utils.GemRelease:
		releaseGem(newVersion, currentVersion, dryRun)
	default:
		releasePlain(newVersion, currentVersion, dryRun)
	}

	utils.Complete()
}
