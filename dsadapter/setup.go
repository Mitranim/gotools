package dsadapter

/********************************** Globals **********************************/

// Config variable. Set on a Setup call.
var conf Config

/*********************************** Setup ***********************************/

// Configures the package. The user calls this once, supplying the config. An
// error is returned if any part of the setup process fails.
func Setup(config Config) error {
	// Save config into global.
	conf = config
	return nil
}
