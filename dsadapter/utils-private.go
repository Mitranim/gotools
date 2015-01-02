package dsadapter

// Private utilities.

import (
	// Standard
	"net/http"
	// Third party
	"github.com/Mitranim/gotools/utils"
)

/********************************* Constants *********************************/

// Error constants.
const (
	err400 = utils.Error("400 bad request")
	err403 = utils.Error("403 insufficient permissions")
	err404 = utils.Error("404 not found")
	err500 = utils.Error("500 internal server error")
)

/********************************* Utilities *********************************/

// Calls the Debugger function from the config, if defined.
func log(req *http.Request, values ...interface{}) {
	if conf.Debugger != nil {
		conf.Debugger(req, values...)
	}
}

// Returns the result of calling the RndId function from the config, if it's
// defined. Otherwise uses the built-in.
func rndId() string {
	if conf.RndId != nil {
		return conf.RndId()
	}
	return RndId()
}

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
