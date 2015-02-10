package render

// Configuration.

import (
	// Standard
	"html/template"
)

/********************************** Config ***********************************/

// Config passed by the user into a Setup call.
type Config struct {
	// Delimiters for templates.
	Delims []string
	// Funcs map for templates.
	Funcs template.FuncMap
	// Directory with hierarchical templates for rendering.
	TemplateDir string
	// Directory with files to read into memory for inlining.
	InlineDir string
	// Function to use for converting integer http status codes to template paths.
	// If omitted, the default CodePath function is used.
	CodePath func(int) string
	// Logging function to use on 500 errors. Pass config.Log to use the default,
	// which is effectively a wrapper for println(). Only works on runtimes that
	// support logging to stdout.
	Logger func(...interface{})
	// Function to check if we're in a development environment. This is checked on
	// each inline call. If true, the file to be inlined is re-read from the disk.
	DevChecker func() bool
	// Bytes to send when rendering fails completely and a hard-set message needs
	// to be written. If omitted, the default err500ISE is used (see `utils-
	// private.go`).
	UltimateFailure []byte
}

/*********************************** Setup ***********************************/

// Generates a state object encapsulating stateful data like parsed templates
// and options. The user calls this once and keeps a reference to the resulting
// state object. An error is returned if any part of the setup process fails.
func Setup(config Config) (State, error) {
	// Create a state object to encapsulate the configuration.
	state := &stateInstance{
		temps:  template.New(""),
		files:  map[string][]byte{},
		config: config,
	}

	// Set up delimiters.
	if len(config.Delims) == 2 {
		state.temps.Delims(config.Delims[0], config.Delims[1])
	}

	// Set up default funcs.
	funcs := makeTemplateFuncs(state)
	state.temps.Funcs(funcs)

	// Set up user funcs.
	if config.Funcs != nil {
		state.temps.Funcs(config.Funcs)
	}

	// Read templates.
	if config.TemplateDir != "" {
		if err := readTemplates(config.TemplateDir, state.temps); err != nil {
			return nil, err
		}
	}

	// Read inline files.
	if config.InlineDir != "" {
		if err := readInline(config.InlineDir, state.files); err != nil {
			return nil, err
		}
	}

	// Return the state object.
	return state, nil
}
