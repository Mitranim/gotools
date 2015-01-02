package render

// Public utilities and constants.

import (
	// Third party
	"github.com/Mitranim/gotools/utils"
)

/********************************* Functions *********************************/

// Default logging function (only for runtimes that support logging to stdout).
var Log = utils.Log

// Converts an error to an http status code.
var ErrorCode = utils.ErrorCode

// Combines the status code and code path functions to generate a template path
// from an error.
func ErrorPath(err error) string {
	return codePath(ErrorCode(err))
}

// Converts an http status code to a template path.
var CodePath = utils.CodePath
