package dsadapter

// State utilities.

import (
	// Standard
	"net/http"
	"reflect"
)

// If the given value is non-nil and has a computer interface, this calls its
// Compute method. If the given value is a slice of computers, this calls the
// Compute method on each element.
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

// Converts the given slice into a slice of records. If the value is not a
// slice, this returns nil. Non-Record elements are ignored. The result is not
// guaranteed to have the same length as the original.
func (this *StateInstance) ToRecords(value interface{}) []Record {
	val := reflect.ValueOf(value)

	// Deference pointer, if any.
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Make sure we're dealing with a slice.
	if val.Kind() != reflect.Slice {
		return nil
	}

	// Make a new value to hold the records.
	records := make([]Record, 0, val.Len())

	// Loop over the old slice and copy Records.
	for i := 0; i < val.Len(); i++ {
		record := val.Index(i).Interface().(Record)

		// Ignore a non-Record.
		if record == nil {
			continue
		}

		// Append the Record.
		records = append(records, record)
	}

	return records
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
	}
	Log(req, values...)
}
