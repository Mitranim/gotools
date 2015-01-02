package dsadapter

// Public utilities and constants.

import (
	// Standard
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"time"

	// App Engine
	"appengine"

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

// Default id function. Makes a random unique id.
func RndId() string {
	i := rand.Int()
	return strconv.Itoa(i)
}

// Republish the error-to-code converter.
var ErrorCode = utils.ErrorCode

// Logs to a GAE context with level `debug`.
func Log(req *http.Request, values ...interface{}) {
	gc := appengine.NewContext(req)
	gc.Debugf(repeat("%v", len(values)), values...)
}
