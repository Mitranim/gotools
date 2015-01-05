package dsadapter

// Configuration.

import (
	// Standard
	"net/http"
)

/********************************** Config ***********************************/

// Config passed by the user into a Setup call.
type Config struct {
	// Function to call when generating a missing id for a new record. If omitted,
	// the default dsadapter.RndId function is used. To disable automatic id generation
	// (not recommended), pass a function that returns an empty string.
	RndId func() string
	// Logger function to call on populate and critical errors. If omitted, no
	// logging is done. Pass dsadapter.Log to use the default (recommended).
	Logger func(*http.Request, ...interface{})
}

/*********************************** Setup ***********************************/

// Generates an object encapsulating stateful data. The user calls this once and
// keeps a reference to the returned object. An error is returned if any part of
// the setup process fails. The state object's methods comprise most of the
// public API of the package.
func Setup(config Config) State {
	return &stateInstance{
		resources:     map[string]Record{},
		populateFuncs: map[string]func(*http.Request){},
		config:        config,
	}
}
