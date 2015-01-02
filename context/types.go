package context

import (
	// Standard
	"encoding/json"
	"net/http"
	"reflect"
)

/********************************** Context **********************************/

// A type that implements the Context interface.
type ContextInstance struct {
	data map[string]interface{}
	rw   http.ResponseWriter
	req  *http.Request
}

/****************************** Context Methods ******************************/

/**
 * Stored values
 */

// Returns the data associated with the request.
func (this *ContextInstance) Data() map[string]interface{} {
	return this.data
}

// Returns the http.ResponseWriter associated with the request.
func (this *ContextInstance) RW() http.ResponseWriter {
	return this.rw
}

// Returns the *http.Request associated with the request.
func (this *ContextInstance) Req() *http.Request {
	return this.req
}

/**
 * Writing
 */

// Sets the http status code to the given integer.
func (this *ContextInstance) Code(code int) Context {
	this.rw.WriteHeader(code)
	return Context(this)
}

// Writes the given string directly to own response writer.
func (this *ContextInstance) Send(str string) Context {
	this.rw.Write([]byte(str))
	return Context(this)
}

// Calls the response writer's Write method with zero input. This marks it as
// having been written to, and causes Martini to terminate the request handler
// chain.
func (this *ContextInstance) End() Context {
	this.rw.Write([]byte{})
	return Context(this)
}

// Renders the template at the given path, writing the output to the
// http.ResponseWriter associated with the current request.
func (this *ContextInstance) Render(path string) {
	bytes, err := render(path, this.Data())
	this.Code(ErrorCode(err))
	this.RW().Write(bytes)
}

/**
 * Error handling
 */

// Sets the status code corresponding to the error and sends its message.
func (this *ContextInstance) SendError(err error) {
	this.Code(ErrorCode(err))
	this.Send(err.Error())
}

// Handles an error. When called with a non-nil error, renders the matching
// error page and sets the status code as per RenderError, then panics with the
// hardcoded message to terminate the caller routine, which should recover with
// `defer Recover()`. When called with a nil error, this is a no-op.
func (this *ContextInstance) Must(err error) {
	// No error -> no-op.
	if err == nil {
		return
	}
	// Error -> render appropriate page, then panic.
	this.Render(codePath(ErrorCode(err)))
	Panic()
}

// Version of Must that sends the error text instead of rendering an error page.
// Intended for JSON API.
func (this *ContextInstance) Ought(err error) {
	// No error -> no-op.
	if err == nil {
		return
	}

	// Error -> set status and write error message, then panic.
	code := ErrorCode(err)
	bytes := []byte(err.Error())
	this.Code(code).RW().Write(bytes)

	Panic()
}

/**
 * HTTP
 */

func (this *ContextInstance) Redirect(path string) {
	http.Redirect(this.rw, this.req, path, http.StatusFound)
}

func (this *ContextInstance) RedirectPermanent(path string) {
	http.Redirect(this.rw, this.req, path, http.StatusMovedPermanently)
}

/**
 * JSON
 */

// Sends the given value as json. If the value is nil, sends a placeholder value
// obtained by checking the value's type with reflection. If decoding fails,
// ends the request with 500 and an empty response.
func (this *ContextInstance) SendAsJson(value interface{}) {
	// Try to encode and fail with 500 if can't.
	bytes, err := json.Marshal(value)
	if err != nil {
		this.Code(500).End()
		return
	}

	val := reflect.ValueOf(value)

	// Dereference the pointer, if any.
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// If value is nil, json.Marshal writes null. If the value has a slice type,
	// we want to write an empty array instead, and for structs and maps, we want
	// to write an empty hash.
	if isNil(value) {
		switch val.Kind() {
		case reflect.Slice, reflect.Array:
			bytes = []byte("[]")
		case reflect.Struct, reflect.Map:
			bytes = []byte("{}")
		}
	}

	this.rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	this.rw.Write(bytes)
}

// Parses the request body as json and writes it to the given destination, which
// must be a pointer. If the json is malformed or writing fails, this
// automatically sends an empty 400 response and ends the request. This does
// have the potential for server errors (non-writable destination) to be
// confused with malformed json errors. ToDo figure out a fix.
func (this *ContextInstance) ParseAsJson(dst interface{}) {
	decoder := json.NewDecoder(this.req.Body)
	err := decoder.Decode(dst)
	if err != nil {
		this.Code(400).End()
	}
}
