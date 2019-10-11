/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package version

import (
	"fmt"

	"github.com/ShogunPanda/impacca/utils"
	"github.com/ShogunPanda/tempera"
	"github.com/spf13/cobra"
)

// InitCLI initializes the CLI
func InitCLI() *cobra.Command {
	cmd := &cobra.Command{
		Use: "version [version]", Aliases: []string{"v"}, Short: "Show or set the current version.", Args: cobra.MaximumNArgs(1), Run: manageVersion,
	}

	cmd.AddCommand(&cobra.Command{Use: "list", Aliases: []string{"a", "all", "l"}, Short: "Show all versions.", Run: listVersion})
	cmd.AddCommand(&cobra.Command{Use: "raw", Aliases: []string{"r"}, Short: "Only show the raw version number.", Run: showRawVersion})

	return cmd
}

func showRawVersion(cmd *cobra.Command, args []string) {
	fmt.Println(utils.GetCurrentVersion())
}

func listVersion(cmd *cobra.Command, args []string) {
	versions := utils.GetVersions()

	utils.Info("Found {secondary}%d{-} versions(s):", len(versions))

	for _, version := range versions {
		fmt.Printf(tempera.ColorizeTemplate("\u0020\u0020\u0020* {primary}%s{-}\n"), version)
	}
}

func manageVersion(cmd *cobra.Command, args []string) {
	currentVersion := utils.GetCurrentVersion()
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	utils.Info("Current version is: {primary}%s{-}", currentVersion)

	if len(args) == 0 {
		return
	}

	if !dryRun {
		utils.GitMustBeClean("change the version")
	}

	newVersion := utils.ChangeVersion(currentVersion, args[0])

	utils.UpdateVersion(newVersion, currentVersion, dryRun)
	utils.Complete()
}
