package context

// Public utility functions and constants.

import (
	// Standard
	"net/http"
	// Third party
	"github.com/Mitranim/gotools/utils"
)

/********************************** Config ***********************************/

// A Config must be passed into each NewContext call. It becomes associated with
// the new context. It lets the user define how the context renders pages.
// context.
type Config struct {
	// Rendering function to use for rendering pages by paths. If omitted, zero
	// bytes are written on each Render call.
	Render func(string, map[string]interface{}) ([]byte, error)

	// Function to convert http error codes into template paths. If omitted, the
	// default CodePath function is used for straight int-to-string conversion.
	CodePath func(int) string

	// Logging function to use on 500 errors. Pass config.Log to use the default,
	// which is effectively a wrapper for println(). Only works on runtimes that
	// support logging to stdout.
	Logger func(...interface{})
}

/********************************* Utilities *********************************/

// Creates a new context from the given http objects.
func NewContext(rw http.ResponseWriter, req *http.Request, config Config) Context {
	return &ContextInstance{
		data:   map[string]interface{}{},
		rw:     rw,
		req:    req,
		config: config,
	}
}

// Recovers from a panic and checks its message. If the panic was intended,
// leaves the application in the recovered state. If the panic was not intended,
// re-panics with the same message.
func Recover() {
	msg := recover()
	if msg != nil && msg != intentionalPanicMessage {
		panic(msg)
	}
}

// Panics with the intentional panic message.
func Panic() {
	panic(intentionalPanicMessage)
}

// Default logging function (only for runtimes that support logging to stdout).
var Log = utils.Log

// Converts an error to an http status code.
var ErrorCode = utils.ErrorCode
