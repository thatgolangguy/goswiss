/*
 * Copyright (c) 2025 Anurag Yadav <contact@anuragyadav.in>
 * License: MIT
 *
 * This file is part of the goswiss project. See LICENSE for details.
 */
 
 package readerutils

import (
	"bufio"
	"errors"
	"os"
)

// The maximum size of the line, if you need to read a really large line then increase this value
const MaxLineSize = 10 * 1024 * 1024 // 10 MB per line

// LineByLineReader reads a text file line by line and returns a slice of strings,
// each representing one line in the file.
func LineByLineReader(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Set the max buffer capacity
	buf := make([]byte, MaxLineSize)
	scanner.Buffer(buf, MaxLineSize)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// StreamLines reads a text file line by line and streams each line into the provided channel
// without loading the entire file into memory. It is designed to handle very large files with
// potentially very long lines by increasing the internal buffer capacity.
//
// The caller must provide a receiving channel of type `chan string`. The channel will be closed
// automatically once the file has been completely read or if an error occurs internally.
//
// Example usage:
//
//	lines := make(chan string)
//	go fileutils.StreamLines("huge-file.txt", lines)
//
//	for line := range lines {
//	    fmt.Println(line)
//	}
//
// StreamLines reads a file line by line and sends each line to a channel without
// loading the entire file into memory. Designed for very large files and long lines.
func StreamLines(filePath string, out chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	// Set the max buffer capacity
	buf := make([]byte, MaxLineSize)
	scanner.Buffer(buf, MaxLineSize)

	defer close(out)
	defer file.Close()

	for scanner.Scan() {
		out <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		panic(errors.New("scanner failed: " + err.Error()))
	}

	return nil
}
