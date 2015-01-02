package dsadapter

// Public utilities and constants.

import (
	// Standard
	"math/rand"
	"reflect"
	"strconv"
	"time"

	// Third party
	"github.com/Mitranim/gotools/utils"
)

// Weak-entropy seed.
func init() {
	rand.Seed(time.Now().UnixNano())
}

/********************************* Constants *********************************/

// CRUD operation codes.
const (
	CodeCreate = iota
	CodeRead
	CodeUpdate
	CodeDelete
)

/********************************* Utilities *********************************/

// Default id function. Makes a random unique id.
func RndId() string {
	i := rand.Int()
	return strconv.Itoa(i)
}

// If the given value is non-nil and has a computer interface, this calls its
// Compute method. If the given value is a slice of computers, this calls the
// Compute method on each element.
func Compute(value interface{}) {
	if value == nil {
		return
	}

	// If the value is a computer, call the method and quit.
	if comp, ok := value.(computer); ok {
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
		Compute(val.Index(i).Interface())
	}
}

// Converts the given slice into a slice of records. If the value is not a
// slice, this returns nil. Non-Record elements are ignored. The result is not
// guaranteed to have the same length as the original.
func ToRecords(value interface{}) []Record {
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

// Republish the error-to-code converter.
var ErrorCode = utils.ErrorCode

// Republish a simple logging function so it can be passed in a Config.
var Log = utils.Log
