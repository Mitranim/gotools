package context

import (
	// Standard
	// "encoding/json"
	"net/http"
)

/********************************** Context **********************************/

// Context interface definition.
type Context interface {

	/**
	 * Stored values
	 */

	// Returns the data associated with the context.
	Data() map[string]interface{}

	// Returns the response writer.
	RW() http.ResponseWriter

	// Returns the request.
	Req() *http.Request

	/**
	 * Writing
	 */

	// Side effect: must set the specified status code on own response writer.
	Code(int) Context

	// Side effect: must write the given string to own response writer.
	Send(string) Context

	// Side effect: must call the response writer's Write method with zero data to
	// mark it as having been written to. This causes Martini to terminate the
	// handler chain after the current handler exits.
	End() Context

	// Takes a path to a page. Must render the page at that path to own
	// http.ResponseWriter. If rendering fails, must render the appropriate error
	// page and set the code as per RenderError.
	Render(string)

	/**
	 * Error handling
	 */

	// Side effect: must set the http status code corresponding to the error type
	// and send the error's message as plain text. Intended for use in API
	// handlers.
	SendError(error)

	// When called with a non-nil error, this must render the error page and set
	// the status code as per RenderError, and panic with the hardcoded message to
	// terminate the caller routine, which should recover with `defer Recover()`.
	// When called with a nil error, this must be a no-op.
	Must(error)

	// Version of Must that, instead of rendering the matching error page, sends
	// the plain error message with the matching status code, as per SendError.
	// When called with a nil error, this must be a no-op. Intended for use in
	// API handlers.
	Ought(error)

	/**
	 * HTTP
	 */

	// Side effect: must redirect to the given path with the code 301.
	Redirect(string)

	// Side effect: must redirect to the given path with the code 302.
	RedirectPermanent(string)

	/**
	 * JSON
	 */

	// Side effect: must write the given value as json and set the Content-Type
	// header to "application/json; charset=UTF-8". If encoding fails, must send
	// error 500 instead.
	SendAsJson(interface{})

	// Side effect: must decode the body of the current request as json and write
	// it to the given destination value, if it's writable. If decoding fails,
	// must end the request with error 400.
	ParseAsJson(interface{})
}
