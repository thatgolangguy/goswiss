/*
 * Copyright (c) 2025 Anurag Yadav <contact@anuragyadav.in>
 * License: MIT
 *
 * This file is part of the goswiss project. See LICENSE for details.
 */

package fileutils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

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
	stream := NewJSONStream()
	var (
		err      error
		template T
	)

	go func() {
		for data := range stream.Watch() {
			if data.Error != nil {
				log.Printf("Error: %s", data.Error.Error())
				err = data.Error
				return
			}
			records = append(records, data.JSONData.(T))
		}
	}()

	stream.Start(filePath, template)
	return records, err
}

// Entry represents each stream. If the stream fails, an error will be present.
type Entry struct {
	Error    error
	JSONData any
}

// Stream helps transmit each streams withing a channel.
type Stream struct {
	stream chan Entry
}

// NewJSONStream returns a new `Stream` type.
func NewJSONStream() Stream {
	return Stream{
		stream: make(chan Entry),
	}
}

// Watch watches JSON streams. Each stream entry will either have an error or a
// User object. Client code does not need to explicitly exit after catching an
// error as the `Start` method will close the channel automatically.
func (s Stream) Watch() <-chan Entry {
	return s.stream
}

// Start starts streaming JSON file line by line. If an error occurs, the channel
// will be closed.
func (s Stream) Start(path string, template any) {
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
