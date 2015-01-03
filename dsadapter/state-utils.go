package dsadapter

// State utilities.

import (
	// Standard
	"net/http"
	"reflect"
)

// If the given value is non-nil and has a computer interface, this calls its
// Compute method. If the given value is a slice of computers, this calls each
// element's Compute method.
func (this *StateInstance) Compute(value interface{}) {
	if value == nil {
		return
	}

	// If the value is a computer, call the method and quit.
	if comp, ok := value.(interface {
		Compute()
	}); ok {
		comp.Compute()
		return
	}

	// Continue to see if this is a slice of computers.
	val := reflect.ValueOf(value)

	// "Dereference" the pointer, if any.
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ignore non-slices.
	if val.Kind() != reflect.Slice {
		return
	}

	// Call Compute on each element.
	for i := 0; i < val.Len(); i++ {
		this.Compute(val.Index(i).Interface())
	}
}

// Returns the result of calling the RndId function from the config, if it's
// defined. Otherwise uses the built-in generator.
func (this *StateInstance) RndId() string {
	if this.config.RndId != nil {
		return this.config.RndId()
	}
	return RndId()
}

/*--------------------------------- Private ---------------------------------*/

// Logs to a GAE context with level `debug`.
func (this *StateInstance) log(req *http.Request, values ...interface{}) {
	if this.config.Logger != nil {
		this.config.Logger(req, values...)
		return
	}
	Log(req, values...)
}
