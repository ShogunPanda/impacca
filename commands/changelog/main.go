/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package changelog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ShogunPanda/impacca/utils"
	"github.com/ShogunPanda/tempera"
	"github.com/spf13/cobra"
)

// InitCLI initializes the CLI
func InitCLI() *cobra.Command {
	cmd := &cobra.Command{Use: "changelog", Aliases: []string{"c"}, Short: "Manage changelog entries.", Run: showChanges}

	cmd.AddCommand(&cobra.Command{
		Use: "list", Aliases: []string{"l"}, Short: "List changelog changes since last version.",
		Run: showChanges,
	})

	cmd.AddCommand(&cobra.Command{
		Use: "version <version>", Aliases: []string{"v"}, Short: "List changelog changes in a specific version.",
		Args: cobra.ExactArgs(1), Run: showVersion,
	})

	cmd.AddCommand(&cobra.Command{
		Use: "save <version> [changes...]", Aliases: []string{"s"}, Short: "Insert all changes since the last version in the CHANGELOG.md file.",
		Args: cobra.MinimumNArgs(1), Run: saveChanges,
	})

	cmd.AddCommand(&cobra.Command{
		Use: "regenerate", Aliases: []string{"r"}, Short: "Regenerates the entire CHANGELOG.md file, EXCLUDING changes since the last version.",
		Run: regenerate,
	})

	return cmd
}

func showChanges(cmd *cobra.Command, args []string) {
	currentVersion := utils.GetCurrentVersion()
	changes := utils.ListChanges(currentVersion.String(), "")
	utils.Info("Found {secondary}%d{-} change(s) since release {secondary}%s{-}:", len(changes), currentVersion)

	for _, change := range changes {
		fmt.Printf(tempera.ColorizeTemplate("\u0020\u0020\u0020* {gray}%s{-}: {primary}%s{-} ({secondary}%s{-})\n"), change.Type, change.Message, change.Hash)
	}
}

func showVersion(cmd *cobra.Command, args []string) {
	currentVersion := args[0]
	versions := utils.GetVersions()

	currentIndex := -1
	for i, v := range versions {
		if v.String() == currentVersion {
			currentIndex = i
			break
		}
	}

	var changes []utils.Change

	if currentIndex > 0 {
		previousVersion := versions[currentIndex-1]

		changes = utils.ListChanges(currentVersion, previousVersion.String())
		utils.Info(
			"Found {secondary}%d{-} change(s) between release {secondary}%s{-} and {secondary}%s{-}:",
			len(changes), previousVersion, currentVersion,
		)
	} else {
		changes = utils.ListChanges(currentVersion, utils.GetFirstCommitHash())

		utils.Info(
			"Found {secondary}%d{-} change(s) between the beginning and release {secondary}%s{-}:",
			len(changes), currentVersion,
		)
	}

	for _, change := range changes {
		fmt.Printf(tempera.ColorizeTemplate("\u0020\u0020\u0020* {gray}%s{-}: {primary}%s{-} ({secondary}%s{-})\n"), change.Type, change.Message, change.Hash)
	}
}

func saveChanges(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	rawChanges := args[1:]
	currentVersion := utils.GetCurrentVersion()
	newVersion := utils.ChangeVersion(currentVersion, args[0])

	changes := make([]utils.Change, 0)

	if len(rawChanges) == 0 {
		changes = utils.ListChanges(currentVersion.String(), "")
	} else {
		for _, c := range rawChanges {
			changes = append(changes, utils.Change{Hash: "", Message: c})
		}
	}

	if !dryRun {
		utils.GitMustBeClean("upload CHANGELOG.md file")
	}

	utils.SaveChanges(newVersion, currentVersion, changes, dryRun)
	utils.Complete()
}

func regenerate(cmd *cobra.Command, args []string) {
	cwd, _ := os.Getwd()
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	versions := utils.GetVersions()

	changelog := ""
	for i, version := range versions {
		previousVersion := ""

		if i > 0 {
			previousVersion = versions[i-1].String()
		} else {
			previousVersion = utils.GetFirstCommitHash()
		}

		date := utils.GetVersionDate(version)
		changes := utils.ListChanges(version.String(), previousVersion)
		changelog = utils.FormatChanges(changelog, version, changes, date)
	}

	if utils.NotifyExecution(dryRun, "Will rewrite", "Rewriting", " CHANGELOG.md file ...") {
		// Save the new file
		err := ioutil.WriteFile(filepath.Join(cwd, "CHANGELOG.md"), []byte(changelog), 0644)

		if err != nil {
			utils.Fatal("Cannot update file {errorPrimary}CHANGELOG.md{-}: {errorPrimary}%s{-}", err.Error())
		}
	}
}
