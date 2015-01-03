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
	// Directory with hierarchical pages for rendering.
	PageDir string
	// Directory with standalone pages for rendering.
	StandaloneDir string
	// Directory with files to read into memory for inlining.
	InlineDir string
	// Function to use for converting integer http status codes to template paths.
	// If omitted, the default CodePath function is used.
	CodePath func(int) string
	// Logging function to use on 500 errors. Pass render.Log to use the default.
	// If omitted, no logging is done. Only works on runtimes that support logging
	// to stdout.
	Logger func(...interface{})
	// Function to check if we're in a development environment. This is checked on
	// each inline call. If true, the file to be inlined is re-read from the disk.
	DevChecker func() bool
	// Bytes to send when rendering fails completely and a hard-set message needs
	// to be written. If omitted, the default err500ISE is used.
	UltimateFailure []byte
}

/*********************************** Setup ***********************************/

// Generates a state object encapsulating stateful data like parsed templates
// and options. The user calls this once and keeps a reference to the resulting
// state object. An error is returned if any part of the setup process fails.
func Setup(config Config) (State, error) {
	// Create a state object to encapsulate the configuration.
	state := &StateInstance{
		pages:       template.New(""),
		standalone:  template.New(""),
		inlineFiles: map[string]template.HTML{},
		config:      config,
	}

	// Set up delimiters.
	if len(config.Delims) == 2 {
		state.pages.Delims(config.Delims[0], config.Delims[1])
		state.standalone.Delims(config.Delims[0], config.Delims[1])
	}

	// Set up default funcs.
	funcs := makeTemplateFuncs(state)
	state.pages.Funcs(funcs)
	state.standalone.Funcs(funcs)

	// Set up user funcs.
	if config.Funcs != nil {
		state.pages.Funcs(config.Funcs)
		state.standalone.Funcs(config.Funcs)
	}

	// Read page templates.
	if config.PageDir != "" {
		if err := readTemplates(config.PageDir, state.pages); err != nil {
			return nil, err
		}
	}

	// Read standalone templates.
	if config.StandaloneDir != "" {
		if err := readTemplates(config.StandaloneDir, state.standalone); err != nil {
			return nil, err
		}
	}

	// Read inline files.
	if config.InlineDir != "" {
		if err := readInline(config.InlineDir, state.inlineFiles); err != nil {
			return nil, err
		}
	}

	// Return the state object.
	return state, nil
}
