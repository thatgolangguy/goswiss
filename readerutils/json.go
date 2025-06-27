/*
 * Copyright (c) 2025 Anurag Yadav <contact@anuragyadav.in>
 * License: MIT
 *
 * This file is part of the goswiss project. See LICENSE for details.
 */

package readerutils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

// Entry represents each stream. If the stream fails, an error will be present.
type Entry struct {
	Error    error
	JSONData any
}

// Stream helps transmit each streams withing a channel.
type Stream struct {
	stream chan Entry
}

// ReadJSONFile reads a JSON array from the specified file path and decodes each element into the generic type T.
// It returns a slice of decoded T objects and any error encountered during parsing.
//
// The JSON file must contain a top-level array of objects (e.g., [{...}, {...}, ...]).
//
// Example usage:
//
//	type Record struct {
//	    ID   int    `json:"id"`
//	    Name string `json:"name"`
//	}
//
//	var records []Record
//	records, err := ReadJSONFile[Record]("data.json", records)
//
// Internally, the function streams the JSON array using a channel and populates
// the provided slice with decoded entries of type T.
//
// NOTE: the records being passed should have exported fields; if they are not then the value is not reflected; record.Name works; record.name won't work.
func ReadJSONFile[T any](filePath string, records []T) ([]T, error) {
	stream := newJSONStream()
	var (
		err      error
		template T
	)

	go func() {
		for data := range stream.watch() {
			if data.Error != nil {
				log.Printf("Error: %s", data.Error.Error())
				err = data.Error
				return
			}
			records = append(records, data.JSONData.(T))
		}
	}()

	stream.start(filePath, template)
	return records, err
}

// StreamJSON reads a JSON array from the specified file path and streams each element into the provided channel.
// The JSON must be a top-level array of objects. The generic type T must have exported fields and match the JSON structure.
//
// Example usage:
//
//	records := make(chan Record)
//	go func() {
//	    for record := range records {
//	        fmt.Println(record)
//	    }
//	}()
//	err := fileutils.StreamJSON("data.json", records)
//	if err != nil {
//	    log.Fatal(err)
//	}
func StreamJSON[T any](filePath string, outputChan chan<- T) error {
	var template T

	// Open file to read
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	defer close(outputChan) // Automatically close user channel when done

	decoder := json.NewDecoder(file)

	// Read opening token (`[`)
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("decode opening delimiter: %w", err)
	}

	i := 1
	for decoder.More() {
		t := reflect.TypeOf(template)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		valuePtr := reflect.New(t).Interface()

		if err := decoder.Decode(valuePtr); err != nil {
			return fmt.Errorf("decode line %d: %w", i, err)
		}

		outputChan <- reflect.ValueOf(valuePtr).Elem().Interface().(T)
		i++
	}

	// Read closing token (`]`)
	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("decode closing delimiter: %w", err)
	}

	return nil
}

// newJSONStream returns a new `Stream` type.
func newJSONStream() Stream {
	return Stream{
		stream: make(chan Entry),
	}
}

// watch func watches JSON streams. Each stream entry will either have an error or a
// User object. Client code does not need to explicitly exit after catching an
// error as the `Start` method will close the channel automatically.
func (s Stream) watch() <-chan Entry {
	return s.stream
}

// start func starts streaming JSON file line by line. If an error occurs, the channel
// will be closed.
func (s Stream) start(path string, template any) {
	// Stop streaming channel as soon as nothing left to read in the file.
	defer close(s.stream)

	// Open file to read.
	file, err := os.Open(path)
	if err != nil {
		s.stream <- Entry{Error: fmt.Errorf("error opening file: %w", err)}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Read opening delimiter. `[` or `{`
	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode opening delimiter: %w", err)}
		return
	}

	// Read file content as long as there is something.
	i := 1
	for decoder.More() {
		// Create a fresh value of the template type
		t := reflect.TypeOf(template)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		valuePtr := reflect.New(t).Interface()

		if err := decoder.Decode(valuePtr); err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode line %d: %w", i, err)}
			return
		}

		// Send back the dereferenced struct
		s.stream <- Entry{JSONData: reflect.ValueOf(valuePtr).Elem().Interface()}
		i++
	}

	// Read closing delimiter. `]` or `}`
	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode closing delimiter: %w", err)}
		return
	}
}
