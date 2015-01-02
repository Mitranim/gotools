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

	// Returns the group of nested pages templates.
	Pages() *template.Template
	// Returns the group of standalone templates.
	Standalone() *template.Template
	// Returns the map of inline files.
	InlineFiles() map[string]template.HTML
	// Returns the configuration object.
	Config() Config

	/*------------------------------- Rendering -------------------------------*/

	// See `render.go`.

	Render(string, map[string]interface{}) ([]byte, error)
	RenderPage(string, map[string]interface{}) ([]byte, error)
	RenderStandalone(string, map[string]interface{}) ([]byte, error)
	RenderError(error, map[string]interface{}) ([]byte, error)
}

/******************************* StateInstance *******************************/

// A type that implements State.
type StateInstance struct {
	pages       *template.Template
	standalone  *template.Template
	inlineFiles map[string]template.HTML
	config      Config
}

/*------------------------------ Stored Values ------------------------------*/

func (this *StateInstance) Pages() *template.Template             { return this.pages }
func (this *StateInstance) Standalone() *template.Template        { return this.standalone }
func (this *StateInstance) InlineFiles() map[string]template.HTML { return this.inlineFiles }
func (this *StateInstance) Config() Config                        { return this.config }

/*--------------------------------- Private ---------------------------------*/

// Adapter methods for generic functions. See their respective documentation.

func (this *StateInstance) log(values ...interface{})  { log(this, values...) }
func (this *StateInstance) errorPath(err error) string { return errorPath(this, err) }
