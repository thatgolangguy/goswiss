/*
 * Copyright (c) 2025 Anurag Yadav <contact@anuragyadav.in>
 * License: MIT
 *
 * This file is part of the goswiss project. See LICENSE for details.
 */

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/thatgolangguy/goswiss/fileutils"
)

// fileutilsCmd represents the fileutils command
var fileutilsCmd = &cobra.Command{
	Use:   "fileutils",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var sizeUtils = &cobra.Command{
	Use:   "size",
	Short: "run the size command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		size, err := fileutils.GetSizeOf(FilePath, fileutils.GB)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Size: %.2f GB\n", size)

	},
}

func init() {
	rootCmd.AddCommand(fileutilsCmd)

	// Add command to the parent
	fileutilsCmd.AddCommand(sizeUtils)

	// Add Flags
	fileutilsCmd.PersistentFlags().StringVarP(&FilePath, "filepath", "f", "", "define the path from where you want to read the file")

}
