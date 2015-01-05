package render

// A state object encapsulates the contextual package state, allowing the
// package to remain stateless. State methods comprise most of the package's
// public API.

import (
	// Standard
	"html/template"
)

/****************************** State Interface ******************************/

// State is an object returned by a Setup call that encapsulates all stateful
// data like parsed templates, inline files, and configuration parameters, and
// comes with adapter methods that use its state in generic render functions.
type State interface {

	/*----------------------------- Stored Values -----------------------------*/

	// Returns the templates group.
	Templates() *template.Template
	// Returns the map of inline files.
	Files() map[string][]byte
	// Returns the configuration object.
	Config() Config

	/*------------------------------- Rendering -------------------------------*/

	// See `render.go`.

	Render(string, map[string]interface{}) ([]byte, error)
	RenderPage(string, map[string]interface{}) ([]byte, error)
	RenderOne(string, map[string]interface{}) ([]byte, error)
	RenderError(error, map[string]interface{}) ([]byte, error)
}

/******************************* stateInstance *******************************/

// A type that implements State.
type stateInstance struct {
	temps  *template.Template
	files  map[string][]byte
	config Config
}

/*------------------------------ Stored Values ------------------------------*/

func (this *stateInstance) Templates() *template.Template { return this.temps }
func (this *stateInstance) Files() map[string][]byte      { return this.files }
func (this *stateInstance) Config() Config                { return this.config }

/*--------------------------------- Private ---------------------------------*/

// Adapter methods for generic functions. See their respective documentation.

func (this *stateInstance) log(values ...interface{})  { log(this, values...) }
func (this *stateInstance) errorPath(err error) string { return errorPath(this, err) }
