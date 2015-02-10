package dsadapter

// Private utilities.

import (
	// Standard
	"fmt"
	"reflect"
	// Third party
	"github.com/Mitranim/gotools/utils"
)

/********************************* Constants *********************************/

// Error constants.
const (
	err403 = utils.Error("403 insufficient permissions")
	err404 = utils.Error("404 not found")
	err422 = utils.Error("422 unprocessable entry")
	err500 = utils.Error("500 internal server error")
)

/********************************* Utilities *********************************/

// Converts an http.Request.URL.Query to a map of params fit for a datastore
// query.
func toParams(query map[string][]string) map[string]string {
	params := map[string]string{}

	for key, param := range query {
		if len(param) == 0 {
			continue
		}
		params[key] = param[0]
	}

	return params
}

// Repeats the given string N times, joined with spaces.
func repeat(str string, count int) (result string) {
	for ; count > 0; count-- {
		if result == "" {
			result = str
		} else if str != "" {
			result += " " + str
		}
	}
	return
}

// Prints expanded values to standard output. For development purposes.
func prn(values ...interface{}) {
	var result string
	for i := 0; i < len(values); i++ {
		if reflect.ValueOf(values[i]).Kind() == reflect.String {
			result += fmt.Sprintf("%v", values[i])
		} else {
			result += fmt.Sprintf("%#v", values[i])
		}
		if i < len(values)-1 {
			result += " "
		}
	}
	fmt.Println(result)
}

// Returns a reflect.Value of the given value. If the value is a pointer,
// returns a reflect.Value of what it references.
func refValue(value interface{}) reflect.Value {
	val := reflect.ValueOf(value)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}
