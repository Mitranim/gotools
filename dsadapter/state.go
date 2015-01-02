package dsadapter

// A state object encapsulates the contextual package state, allowing the
// package to remain stateless. Most of the package's API is tied to a state
// object that the user must generate.

import (
	// Standard
	"net/http"

	// App Engine
	"appengine/datastore"
)

/****************************** State Interface ******************************/

// State is an object returned by a Setup call that encapsulates stateful data
// and provides most of the package's API.
type State interface {

	/*--------------------------- Record Operations ---------------------------*/

	// See `record.go`.

	Key(*http.Request, Record) *datastore.Key
	Read(*http.Request, Record) error
	Save(*http.Request, Record) error
	Delete(*http.Request, Record) error
	FindOne(*http.Request, Record, map[string]string) error

	/*------------------------- Collection Operations -------------------------*/

	// See `collection.go`.

	Find(*http.Request, interface{}, map[string]string, int) error
	FindAll(*http.Request, interface{}, map[string]string) error
	FindByQuery(*http.Request, interface{}) error

	/*------------------------------- Resources -------------------------------*/

	// See `resource.go`.

	Resources() map[string]Record
	NewRecordByResource(string) Record
	NewCollectionByResource(string) interface{}
	SliceOf(interface{}) interface{}
	NewRecordFromCollection(interface{}) (Record, error)

	/*------------------------------- Populate --------------------------------*/

	// See `populate.go`.

	PopulateFuncs() map[string]func(*http.Request)
	RegisterForPopulate(interface{})
	Populate(*http.Request)

	/*------------------------------- Utilities -------------------------------*/

	// See `state-utils.go`.

	Compute(interface{})
	RndId() string
}

/******************************* StateInstance *******************************/

// A type that implements State.
type StateInstance struct {
	resources     map[string]Record
	populateFuncs map[string]func(*http.Request)
	config        Config
}
