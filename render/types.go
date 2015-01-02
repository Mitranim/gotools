package render

import (
	"html/template"
)

/********************************** Config ***********************************/

// Config passed by the user to the Setup call.
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
	// If omitted, no logging is done.
	Logger func(...interface{})
	// Function to check if we're in a development environment. This is checked on
	// each inline call. If true, the file to be inlined is re-read from the disk.
	DevChecker func() bool
	// Bytes to send when rendering fails completely and a hard-set message needs
	// to be written. If omitted, the default err500ISE is used.
	UltimateFailure []byte
}

/******************************** readWriter *********************************/

// Rudimental io.ReadWriter.
type readWriter []byte

func (this *readWriter) Write(bytes []byte) (int, error) {
	*this = append(*this, bytes...)
	return len(*this), nil
}

func (this *readWriter) Read(bytes []byte) (int, error) {
	copy(bytes, *this)
	return len(bytes), nil
}

func (this *readWriter) String() string {
	return string(*this)
}

/********************************* errorStr **********************************/

// Rudimental error type.
type errorStr string

// Error method.
func (this errorStr) Error() string {
	return string(this)
}
