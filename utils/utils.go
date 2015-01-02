package utils

// Utilities shared between gotools packages.

import (
	"fmt"
)

/********************************* Utilities *********************************/

// Takes an error and tries to get an http status code from it by scanning
// digits from the beginning. Example: Error("403 unauthorised") -> 403 int. If
// the error is nil, this returns 200. If the error code can't be found or lies
// outside the standard http error range (400 <= x <= 599), this returns 500.
func ErrorCode(err error) int {
	if err == nil {
		return 200
	}
	code := atoi(err.Error())
	if code < 400 || code > 599 {
		return 500
	}
	return code
}

// Combines the status code and code path functions to generate a template path
// from an error.
func ErrorPath(err error) string {
	return CodePath(ErrorCode(err))
}

// Converts an http status code into a template path through a straight number-
// to-string translation.
func CodePath(code int) string {
	return itoa(code)
}

// print/println wrapper that automatically does fmt %v for each value.
func Log(values ...interface{}) {
	for _, value := range values {
		print(fmt.Sprintf("%v", value) + " ")
	}
	println()
}
