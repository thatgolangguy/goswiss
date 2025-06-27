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
	"time"

	"github.com/spf13/cobra"
	"github.com/thatgolangguy/goswiss/readerutils"
)

type Record struct {
	Name     string  `json:"name"`
	Language string  `json:"language"`
	ID       string  `json:"id"`
	Bio      string  `json:"bio"`
	Version  float32 `json:"version"`
}

// readerutilsCmd represents the fileutils command
var readerutilsCmd = &cobra.Command{
	Use:   "fileutils",
	Short: "cmd to access fileutils",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fileutils called")
	},
}

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "cmd to read json file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var records []Record
		var err error
		startTime := time.Now()

		records, err = readerutils.ReadJSONFile("example.json", records)
		if err != nil {
			log.Fatalf("failed to read json file: %s\n", err)
		}

		elapsed := time.Since(startTime)
		seconds := elapsed.Seconds()
		recordsPerSecond := float64(len(records)) / seconds

		log.Printf("✅ %d records fetched in %s (%.0f records/sec)",
			len(records),
			elapsed.Truncate(time.Millisecond),
			recordsPerSecond,
		)
	},
}

var streamJsonCmd = &cobra.Command{
	Use:   "stream-json",
	Short: "read the json file and stream it through a channel",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		recordsChan := make(chan Record)
		startTime := time.Now()

		go func() {
			err := readerutils.StreamJSON("example.json", recordsChan)
			if err != nil {
				log.Fatalf("failed to read json file: %s\n", err)
			}
		}()

		var count int64
		for range recordsChan {
			// log.Printf("%+v\n", record)
			count++
		}

		elapsed := time.Since(startTime)
		seconds := elapsed.Seconds()
		recordsPerSecond := float64(count) / seconds

		log.Printf("✅ %d records fetched in %s (%.0f records/sec)",
			count,
			elapsed.Truncate(time.Millisecond),
			recordsPerSecond,
		)

	},
}

var streamLinesCmd = &cobra.Command{
	Use:   "stream-lines",
	Short: "read the json file and stream it through a channel",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		recordsChan := make(chan string)
		startTime := time.Now()

		go func() {
			err := readerutils.StreamLines("test.go", recordsChan)
			if err != nil {
				log.Fatalf("failed to read file: %s\n", err)
			}
		}()

		var count int64
		for record := range recordsChan {
			log.Printf("%+v\n", record)
			count++
		}

		elapsed := time.Since(startTime)
		seconds := elapsed.Seconds()
		recordsPerSecond := float64(count) / seconds

		log.Printf("✅ %d records fetched in %s (%.0f records/sec)",
			count,
			elapsed.Truncate(time.Millisecond),
			recordsPerSecond,
		)

	},
}

func init() {
	rootCmd.AddCommand(readerutilsCmd)

	// Add child commands
	readerutilsCmd.AddCommand(jsonCmd)
	readerutilsCmd.AddCommand(streamJsonCmd)
	readerutilsCmd.AddCommand(streamLinesCmd)
}
