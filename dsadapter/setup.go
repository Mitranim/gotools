package dsadapter

import (
	// Standard
	"net/http"
)

/********************************** Globals **********************************/

// Config variable. Set on a Setup call.
var conf Config

/********************************** Config ***********************************/

// The Config object passed by the user when setting up this package.
type Config struct {
	// Function to call when generating a missing id for a new record. If omitted,
	// the default dsadapter.RndId function is used. To disable automatic id generation
	// (not recommended), pass a function that returns an empty string.
	RndId func() string
	// Logger function to call on populate and critical errors. If omitted, no
	// logging is done. Pass dsadapter.Debug to use the default (recommended).
	Debugger func(*http.Request, ...interface{})
}

/*********************************** Setup ***********************************/

// Configures the package. The user calls this once, supplying the config. An
// error is returned if any part of the setup process fails.
func Setup(config Config) error {
	// Save config into global.
	conf = config
	return nil
}
