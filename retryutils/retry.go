/*
 * Copyright (c) 2025 Anurag Yadav <contact@anuragyadav.in>
 * License: MIT
 *
 * This file is part of the goswiss project. See LICENSE for details.
 */

package retryutils

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// Retry attempts to execute a function multiple times with a specified delay between each attempt.
// It uses reflection to call the given function with parameters, and expects the function's last
// return value to be of type `error`. If the error is nil, the operation is considered successful.
// Otherwise, it retries the operation up to `maxRetries` times.
//
// Parameters:
//   - operation: a function to be executed; it must return an error as its last return value.
//   - params: the parameters to pass to the function, followed by:
//   - maxRetries (int): the number of times to retry on failure
//   - delay (time.Duration): the time to wait between retries
//
// Returns:
//   - []reflect.Value: the return values of the function if it succeeds
//   - error: if all retry attempts fail, the last encountered error is returned
//
// Example usage:
//
//	result, err := Retry(myFunc, arg1, arg2, 3, time.Second)
//
// Note:
//   - The function passed in must have an `error` as its last return value.
//   - `maxRetries` and `delay` must be the last two elements in `params`, respectively.
func Retry(operation any, params ...any) ([]reflect.Value, error) {
	opValue := reflect.ValueOf(operation)
	if opValue.Kind() != reflect.Func {
		return nil, errors.New("operation must be a function")
	}

	if len(params) < 2 {
		return nil, errors.New("not enough parameters provided to Retry function")
	}

	maxRetries, ok := params[len(params)-2].(int)
	if !ok {
		return nil, errors.New("maxRetries (second to last param) must be an int")
	}

	delay, ok := params[len(params)-1].(time.Duration)
	if !ok {
		return nil, errors.New("delay (last param) must be a time.Duration")
	}

	opParams := params[:len(params)-2]
	in := make([]reflect.Value, len(opParams))
	for i, param := range opParams {
		in[i] = reflect.ValueOf(param)
	}

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		out := opValue.Call(in)
		lastReturnValue := out[len(out)-1]

		if !lastReturnValue.IsNil() {
			lastErr, ok = lastReturnValue.Interface().(error)
			if !ok {
				return nil, fmt.Errorf("last return value of operation is not an error but a %T", lastReturnValue.Interface())
			}
		} else {
			return out, nil
		}

		if attempt < maxRetries {
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("operation failed after %d attempts: %w", maxRetries, lastErr)
}
