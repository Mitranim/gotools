package dsadapter

// Interfaces defined in this package.

import (
	// Standard
	"net/http"
)

/******************************** Interfaces *********************************/

// Record is an instance of a database model with common parsing, crud and
// validation methods.
type Record interface {

	/*------------------------------- Lifecycle -------------------------------*/

	// Must return a map of record fields to error messages.
	Validate() map[string]string

	// Modifies and adjusts properties. May change existing fields and create new
	// ones. This must be called every time a record or a collection is read from
	// a database or parsed from source data.
	Compute()

	// Answers whether the user associated with this request has the rights to
	// perform the given CRUD operation on this particular value. See the CRUD
	// constants.
	Can(*http.Request, int) bool

	/*--------------------------------- CRUD ----------------------------------*/

	// Side effect: must save self to the Datastore under the Datastore key
	// identified by own id.
	Save(*http.Request) error

	// Side effect: must read a record by own id from the Datastore into self.
	Read(*http.Request) error

	// Side effect: must delete self by id from the Datastore.
	Delete(*http.Request) error

	// ToDo: consider adding a Patch method that would combine Read() and
	// Save() to update an existing record without deleting missing fields.
	// Step 1: create a clone of self with the same id and Read() it.
	// Step 2: for each zero field on self, copy that field from the clone.
	// Step 3: Save() self.
	// End result: the old record has been updated with non-zero-value fields
	// from the new record. Zero value fields on the new record have been ignored.
	//
	// Leaning towards rejecting the idea because it would add complexity and
	// might lead to programmer errors like using Patch() in situations when
	// writing empty fields is actually intended (e.g. set some text field to "").

	/*------------------------------- Utilities -------------------------------*/

	// Returns own id.
	GetId() string

	// Side effect: must set own id to the given string.
	SetId(string)

	// Returns own Datastore kind.
	Kind() string
}

// Can compute properties.
type computer interface {
	Compute()
}

/**
 * Generic collection definition.
 *
 * Some functions in this package require values of a generic collection type
 * that we can't really define in Go code, so here it is:
 *
 * type Collection []*S,
 *
 * where S is any struct type that implements the Record interface. The proper
 * type is verified with reflect, but the callers should still take care to
 * pass matching values.
 *
 * Take note that []Record is NOT the same type and does not satisfy the
 * definition. We need an underlying struct type for Datastore queries.
 *
 * Also take note that this type is not used directly in any function in this
 * package (we use interface{} instead). It's for reference purposes only.
 */

/**
 * Record lifecycle:
 *   * (created from json or form data -> its Validate method is called ->
 *     sequence fails if validation fails) || read from database || created
 *     from mock data
 *   * its Compute method is called to calculate derived properties
 *   * the user does something with it
 *   * the user calls its Save method
 *     * the Validate method is called; the sequence fails if validation fails
 *     * if the record doesn't have an id, a new random id is generated
 *     * record is saved to database
 *   * the user calls its Delete method
 *     * record is deleted from database
 */

/********************************** Config ***********************************/

// The Config object passed by the user when setting up this package.
type Config struct {
	// Function to call when generating a missing id for a new record. If omitted,
	// the default dsadapter.RndId function is used. To disable automatic id generation
	// (not recommended), pass a function that returns an empty string.
	RndId func() string
	// Logger function to call on populate and critical errors. If omitted, no
	// logging is done. Pass dsadapter.Log to use the default logger.
	Logger func(...interface{})
}

/********************************* errorStr **********************************/

// Rudimental error type.
type errorStr string

// Error method.
func (this errorStr) Error() string {
	return string(this)
}
