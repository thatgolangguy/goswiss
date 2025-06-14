/*
 * Copyright (c) 2025 Anurag Yadav <contact@anuragyadav.in>
 * License: MIT
 *
 * This file is part of the goswiss project. See LICENSE for details.
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thatgolangguy/goswiss/sliceutils"
)

// sliceutilsCmd represents the sliceutils command
var sliceutilsCmd = &cobra.Command{
	Use:   "sliceutils",
	Short: "cmd to access sliceutils",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var chunkCmd = &cobra.Command{
	Use:   "chunk",
	Short: "cmd to access sliceutils",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		records := []int{1, 2, 3, 4, 5, 6, 7}
		chunks, _ := sliceutils.CreateChunks(records, 2)
		fmt.Println(chunks)
	},
}

func init() {
	rootCmd.AddCommand(sliceutilsCmd)

	sliceutilsCmd.AddCommand(chunkCmd)
}
