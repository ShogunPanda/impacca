/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package release

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/ShogunPanda/impacca/utils"
	"github.com/ShogunPanda/tempera"
	"github.com/spf13/cobra"
)

// InitCLI initializes the CLI
func InitCLI() *cobra.Command {
	cmd := &cobra.Command{Use: "release", Aliases: []string{"r"}, Short: "Manage GitHub releases.", Run: showReleases}
	cmd.PersistentFlags().StringP("remote", "r", "origin", "The git remote name.")
	cmd.PersistentFlags().StringP("token", "t", "", "The GitHub API token.")

	cmd.AddCommand(&cobra.Command{
		Use: "show <version>", Aliases: []string{"r"}, Short: "Show GitHub release.",
		Args: cobra.ExactArgs(1), Run: showRelease,
	})

	cmd.AddCommand(&cobra.Command{
		Use: "save <version>", Aliases: []string{"s"}, Short: "Updates all changes in version to the GitHub release",
		Args: cobra.MinimumNArgs(1), Run: saveRelease,
	})

	cmd.AddCommand(&cobra.Command{
		Use: "regenerate", Aliases: []string{"a"}, Short: "Regenerates all GitHub releases using local versions.",
		Run: regenerateReleases,
	})

	return cmd
}

func printRelease(release utils.Release) {
	fmt.Printf(tempera.ColorizeTemplate(fmt.Sprintf(
		"\u0020\u0020\u0020* Version {primary}%s{-} ({secondary}%s{-})\n", 
		release.Version.String(), release.Date.Format("2006-01-02"),
	)))

	if release.Body != "" {
		fmt.Println("\n" + utils.Indent(release.Body, "\u0020\u0020\u0020\u0020\u0020"))
	}

	fmt.Println("")
}

func showReleases(cmd *cobra.Command, args []string) {
	remote, _ := cmd.Flags().GetString("remote")
	repository := utils.DetectGithubRepository(remote, false)

	res := utils.GitHubReleaseAPICall(
		"get GitHub releases", "GET", fmt.Sprintf("/repos/%s/releases", repository), 
		"", map[string]string{}, true,
	)

	var releases []utils.Release
	err := res.JSON(&releases)
	
	if err != nil {
		utils.Fatal("Cannot decode JSON response to get GitHub releases: {errorPrimary}%s{-}", err.Error())
	}	

	if len(releases) == 0 {
		utils.Warn("No GitHub releases found.")
		return
	}

	// Sort release by version, descending
	sort.SliceStable(releases, func(i, j int) bool { return releases[i].Version.GreaterThan(releases[j].Version) })

	utils.Info("Found {secondary}%d{-} GitHub release(s):\n", len(releases))
	
	for _, release := range releases {
		printRelease(release)
	}
}

func showRelease(cmd *cobra.Command, args []string) {
	remote, _ := cmd.Flags().GetString("remote")
	repository := utils.DetectGithubRepository(remote, false)
	version, _ := semver.NewVersion(args[0])

	res := utils.GitHubReleaseAPICall(
		"get a GitHub release", "GET", fmt.Sprintf("/repos/%s/releases/tags/v%s", repository, version), 
		"", map[string]string{}, true,
	)

	if res.StatusCode == 404 {
		utils.Fatal("Cannot find GitHub release {errorPrimary}%s{-}.", version.String())
	}

	var release utils.Release
	err := res.JSON(&release)
	
	if err != nil {
		utils.Fatal("Cannot decode JSON response to get a GitHub release: {errorPrimary}%s{-}", err.Error())
	}	

	utils.Info("Found one GitHub release:\n")
	printRelease(release)
}

func saveRelease(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	remote, _ := cmd.Flags().GetString("remote")
	token, _ := cmd.Flags().GetString("token")
	repository := utils.DetectGithubRepository(remote, false)
	version, _ := semver.NewVersion(args[0])

	utils.SaveRelease(version, repository, remote, token, dryRun)
}

func regenerateReleases(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	remote, _ := cmd.Flags().GetString("remote")
	token, _ := cmd.Flags().GetString("token")
	repository := utils.DetectGithubRepository(remote, false)
	versions := utils.GetVersions()

	for _, version := range versions {
		utils.SaveRelease(version, repository, remote, token, dryRun)
	}
}
