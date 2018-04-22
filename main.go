/*
 * This file is part of impacca. Copyright (C) 2013 and above Shogun <shogun@cowtech.it>.
 * Licensed under the MIT license, which can be found at https://choosealicense.com/licenses/mit.
 */

package main

import (
	"github.com/ShogunPanda/tempera"
	"github.com/spf13/cobra"

	"github.com/ShogunPanda/impacca/commands/changelog"
	"github.com/ShogunPanda/impacca/commands/release"
	"github.com/ShogunPanda/impacca/commands/version"
)

func main() {
	tempera.AddCustomStyle("primary", "bold", "blue")
	tempera.AddCustomStyle("secondary", "bold", "yellow")
	tempera.AddCustomStyle("errorPrimary", "bold", "white")

	var rootCmd = &cobra.Command{Use: "impacca", Short: "Package releasing made easy."}
	rootCmd.Version = "1.0.0"
	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "Do not execute write operation, only show them.")

	rootCmd.AddCommand(version.InitCLI())
	rootCmd.AddCommand(changelog.InitCLI())
	rootCmd.AddCommand(release.InitCLI())

	rootCmd.Execute()
}
