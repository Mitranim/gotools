package context

/**
 * ToDo:
 *
 * Review the mapping helper.
 *
 * Consider adding premapping (for e.g. martini.Params).
 *
 * Config should include an option to take a pointer to a nil interface that
 * mimics Context, and use it for mapping. This is useful for republishing it
 * to make handler definitions shorter.
 *
 * Consider how to make it easier to define and map a user-defined superset of
 * Config.
 */

import ()

/********************************** Globals **********************************/

// Config variable. Set on a Setup call.
var conf Config

// Readiness status. Set to true after a successful Setup call.
var ready bool

/********************************** Config ***********************************/

type Config struct {
	// Rendering function used for rendering pages by paths.
	Render func(string, map[string]interface{}) ([]byte, error)

	// Function to convert http error codes into template paths.
	CodePath func(int) string
}

/*********************************** Init ************************************/

func Setup(config Config) error {
	if ready {
		return nil
	}

	// Save into global.
	conf = config

	ready = true
	return nil
}
