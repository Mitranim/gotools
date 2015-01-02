package render

// Configuration.

import (
	"html/template"
)

/********************************** Globals **********************************/

// Page templates.
var Pages = template.New("Pages")

// Standalone templates.
var Standalone = template.New("Standalone")

// Files for inlining.
var InlineFiles = map[string]template.HTML{}

// Config variable. Set on a Setup call.
var conf Config

// Readiness status. Set to true after a successful Setup call.
var ready bool

/******************************* Configuration *******************************/

// Configures the renderer. The user calls this once, supplying the config. An
// error is returned if any part of the setup process fails. Currently, this
// function can only be called once, and ignores subsequent calls. ToDo move
// template definitions here to enable repeat calls.
func Setup(config Config) error {
	if ready {
		return nil
	}

	// Save into global.
	conf = config

	// Set up delimiters.
	if len(config.Delims) == 2 {
		Pages.Delims(config.Delims[0], config.Delims[1])
		Standalone.Delims(config.Delims[0], config.Delims[1])
	}

	// Set up default funcs.
	Pages.Funcs(templateFuncs)
	Standalone.Funcs(templateFuncs)

	// Set up user funcs.
	if config.Funcs != nil {
		Pages.Funcs(config.Funcs)
		Standalone.Funcs(config.Funcs)
	}

	// Read page templates.
	if config.PageDir != "" {
		if err := readTemplates(config.PageDir, Pages); err != nil {
			return err
		}
	}

	// Read standalone templates.
	if config.StandaloneDir != "" {
		if err := readTemplates(config.StandaloneDir, Standalone); err != nil {
			return err
		}
	}

	// Read inline files.
	if config.InlineDir != "" {
		if err := readInline(config.InlineDir); err != nil {
			return err
		}
	}

	ready = true
	return nil
}
