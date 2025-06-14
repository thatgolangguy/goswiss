/*
 * Copyright (c) 2025 Anurag Yadav <contact@anuragyadav.in>
 * License: MIT
 *
 * This file is part of the goswiss project. See LICENSE for details.
 */

package sliceutils

import "fmt"

// CreateChunks splits a slice of type T into smaller chunks of the specified size.
// Returns a slice of slices ([][]T) and an error if input is invalid.
//
// Example:
//
//	records := []int{1, 2, 3, 4, 5, 6, 7}
//	chunks, _ := CreateChunks(records, 3)
//	// chunks => [[1 2 3] [4 5 6] [7]]
func CreateChunks[T any](records []T, chunkSize int) ([][]T, error) {
	if chunkSize <= 0 {
		return nil, fmt.Errorf("chunk size must be greater than zero")
	}

	var chunks [][]T
	for i := 0; i < len(records); i += chunkSize {
		end := i + chunkSize
		if end > len(records) {
			end = len(records)
		}
		chunks = append(chunks, records[i:end])
	}

	return chunks, nil
}
