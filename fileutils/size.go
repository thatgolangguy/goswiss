/*
 * Copyright (c) 2025 Anurag Yadav <contact@anuragyadav.in>
 * License: MIT
 *
 * This file is part of the goswiss project. See LICENSE for details.
 */

package fileutils

import (
	"os"
)

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
	TB = 1 << 40
)

// GetSizeOf returns the size of the file at the specified path in the unit defined.
//
// Example usage:
//
//	size, err := fileutils.GetSizeOf("large.json", MB)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("File size: %d MB\n", size)
func GetSizeOf(filePath string, unit int64) (float64, error) {
	if info, err := os.Stat(filePath); err != nil {
		return 0, err
	} else {
		return float64(info.Size()) / float64(unit), nil
	}
}
