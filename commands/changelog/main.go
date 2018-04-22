/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package changelog

import (
	"fmt"

	"github.com/ShogunPanda/impacca/utils"
	"github.com/ShogunPanda/tempera"
	"github.com/spf13/cobra"
)

// InitCLI initializes the CLI
func InitCLI() *cobra.Command {
	cmd := &cobra.Command{Use: "changelog", Aliases: []string{"c"}, Short: "Manage changelog entries.", Run: showChanges}

	cmd.AddCommand(&cobra.Command{Use: "list", Aliases: []string{"l"}, Short: "List changelog changes since last version.", Run: showChanges})

	saveCommand := &cobra.Command{
		Use: "save <version> [changes...]", Aliases: []string{"s"}, Short: "Insert all changes since the last version in the CHANGELOG.md file.",
		Args: cobra.MinimumNArgs(1), Run: saveChanges,
	}
	cmd.AddCommand(saveCommand)

	return cmd
}

func showChanges(cmd *cobra.Command, args []string) {
	currentVersion := utils.GetCurrentVersion()
	changes := utils.ListChanges(currentVersion.String())
	utils.Info("Found {secondary}%d{-} change(s) since release {secondary}%s{-}:", len(changes), currentVersion)

	for _, change := range changes {
		fmt.Printf(tempera.ColorizeTemplate("\u0020\u0020\u0020* {primary}%s{-} ({secondary}%s{-})\n"), change.Message, change.Hash)
	}
}

func saveChanges(cmd *cobra.Command, args []string) {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
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
		utils.GitMustBeClean("upload CHANGELOG.md file")
	}

	utils.SaveChanges(newVersion, currentVersion, changes, dryRun)
	utils.Complete()
}
