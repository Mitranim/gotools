package context

// Private utility functions and constants.

import (
	// Standard
	"reflect"
	// Third party
	"github.com/Mitranim/gotools/utils"
)

/********************************* Constants *********************************/

// Message included into intentional panics used to exit the http handler early.
// The Recover function checks the message and repanics if it doesn't match this
// constant.
const intentionalPanicMessage = "Panic() was called to terminate the caller routine; precede with `defer Recover()` to continue app execution"

/********************************* Utilities *********************************/

// Checks if the given value is nil.
func isNil(value interface{}) bool {
	val := reflect.ValueOf(value)

	// Indirect a pointer.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Only call IsNil on supported types.
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

// Render wrapper. Calls the config.Render function if provided, otherwise
// returns zero bytes.
func render(ct *ContextInstance, path string, data map[string]interface{}) ([]byte, error) {
	if ct.config.Render != nil {
		return ct.config.Render(path, data)
	}
	return []byte{}, nil
}

// CodePath wrapper. Calls the config.CodePath function if provided, otherwise
// uses a straight number-to-string conversion.
func codePath(ct *ContextInstance, code int) string {
	if ct.config.CodePath != nil {
		return ct.config.CodePath(code)
	}
	return utils.CodePath(code)
}
