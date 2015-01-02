package context

// Public utility functions and constants.

import (
	// Standard
	"net/http"
	// Third party
	"github.com/Mitranim/gotools/utils"
)

/********************************* Utilities *********************************/

// Creates a new context from the given http objects.
func NewContext(rw http.ResponseWriter, req *http.Request) Context {
	return &ContextInstance{
		data: map[string]interface{}{},
		rw:   rw,
		req:  req,
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

// Panics with the intended panic message.
func Panic() {
	panic(intentionalPanicMessage)
}

// Converts an error to an http status code.
var ErrorCode = utils.ErrorCode
