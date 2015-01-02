package dsadapter

// Private utilities.

import (
	// Third party
	"github.com/Mitranim/gotools/utils"
)

/********************************* Constants *********************************/

// Error constants.
const (
	err403 = utils.Error("403 insufficient permissions")
	err404 = utils.Error("404 not found")
	err422 = utils.Error("422 unprocessable entry")
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
