package render

// Public render methods.

import (
	// Standard
	"html/template"
	"strings"
)

/**
 * Shorthand rendering method. Renders the page at the given path,
 * automatically falling back to error pages corresponding to the kinds of
 * errors that may occur (404, 500, possibly others). Returns the rendered
 * bytes and the last error that occurred in the process.
 *
 * The caller routine should always examine the error with ErrorCode(err) to
 * retrieve the http status code to set in the response handler.
 *
 * Note that rendering is always going to be successful; the role of the error
 * is not to signal a complete failure, but to carry the information about the
 * character of the problem (if any) that occurred in the process.
 *
 * Also see the renderError comment.
 */
func (this *stateInstance) Render(path string, data map[string]interface{}) ([]byte, error) {
	bytes, err := this.RenderPage(path, data)

	if err != nil {
		return this.RenderError(err, data)
	}

	return bytes, nil
}

// Takes a path to a page and a data map. Renders the page and, hierarchically,
// all layouts enclosing it, up to the root, passing the data map to each
// template.
func (this *stateInstance) RenderPage(path string, data map[string]interface{}) ([]byte, error) {
	// Adjust and validate path.
	path, err := parsePath(this.temps, path)
	if err != nil {
		return nil, err
	}

	// Check for nil map.
	if data == nil {
		data = map[string]interface{}{}
	}

	// Build an array of nested template paths.
	paths := pathsToTemplates(path)

	// Render the template into each enclosing layout.
	for _, pt := range paths {
		bytes, err := renderAt(this.temps, pt, data)
		if err != nil {
			return nil, err
		}
		// Enclose the content.
		data["content"] = template.HTML(strings.TrimSpace(string(bytes)))
	}

	// Grab the resulting content as bytes.
	html, _ := data["content"].(template.HTML)

	return []byte(html), nil
}

// Renders a template at the given path, ignoring the page hierarchy.
func (this *stateInstance) RenderOne(path string, data map[string]interface{}) ([]byte, error) {
	return renderAt(this.temps, path, data)
}

/**
 * Takes an error and renders a page at the path corresponding to the error,
 * automatically falling back to other error pages if a different error occurs.
 * Returns the rendered bytes and the last error that occurred in the process.
 *
 * Error codes are translated into template paths by using either the CodePath
 * function provided during setup, or a simple integer-to-string conversion
 * (default). If your error pages are located in the root of the pages folder
 * and have names like "404.html" or "500.gohtml", they will be used
 * automatically.
 *
 * The caller routine should always examine the error with ErrorCode(err) to
 * retrieve the http status code to set in the response handler.
 *
 * Note that rendering is always going to be successful; the role of the error
 * is not to signal a complete failure, but to carry the information about the
 * character of the problem (if any) that occurred in the process.
 */
func (this *stateInstance) RenderError(err error, data map[string]interface{}) (bytes []byte, lastErr error) {
	// Map of error codes that have occurred at least once.
	codes := map[int]bool{}

	/**
	 * Algorithm:
	 *  * attempt to render each non-500 error once; if the same code repeats,
	 *  	fall through to 500
	 *  * attempt to render 500 once; if 500 repeats, fall back on bytes
	 *  * if something is rendered without an error, immediately break and return
	 *  	the result plus the last non-nil error
	 */

	for code := ErrorCode(err); err != nil && !codes[code]; codes[code] = true {
		lastErr = err
		// Try to render the matching page.
		bytes, err = this.RenderPage(this.errorPath(err), data)
	}

	if err == nil {
		return
	}

	// At this point, any error is considered internal.
	lastErr = err500ISE

	// If 500 hasn't occurred yet, try to render it.
	if !codes[500] {
		bytes, err = this.RenderPage(this.errorPath(err), data)
	}

	if err == nil {
		return
	}

	// If rendering of page 500 fails, we fall back on bytes.
	this.log("internal rendering error:", err)
	// Use the provided UltimateFailure data, if possible.
	if len(this.config.UltimateFailure) > 0 {
		bytes = this.config.UltimateFailure
		// Otherwise use the default message.
	} else {
		bytes = []byte(err500ISE)
	}

	return
}
