package dsadapter

// Private utilities.

import ()

/********************************* Constants *********************************/

// Error constants.
const (
	err400 = errorStr("400 bad request")
	err403 = errorStr("403 insufficient permissions")
	err404 = errorStr("404 not found")
	err500 = errorStr("500 internal server error")
)

/********************************* Utilities *********************************/

// Calls the Logger function from the config, if it's defined.
func log(values ...interface{}) {
	if conf.Logger != nil {
		conf.Logger(values...)
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
