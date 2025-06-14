package fileutils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

// ReadJSONUsingChan reads a JSON array from the specified file path and streams each element into the provided channel.
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
//	err := fileutils.ReadJSONUsingChan("data.json", records)
//
//	if err != nil {
//	    log.Fatal(err)
//	}
func ReadJSONUsingChan[T any](filePath string, outputChan chan<- T) error {
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
