package dsadapter

// Record related functions.

import (
	// Standard
	"net/http"
	"reflect"

	// App Engine
	"appengine"
	"appengine/datastore"
)

/****************************** Query Functions ******************************/

// Takes a pointer to a record and tries to find one record of the matching
// type, filtered by the given params. If a record is successfully found, it's
// written to the destination, which must be a pointer. If not, an error if
// returned.
func FindOne(req *http.Request, destination Record, params map[string]string) error {
	// Make a matching collection.
	collection := SliceOf(destination)

	// Try to find one of that type.
	err := Find(req, collection, params, 1)
	if err != nil {
		return err
	}

	// Prepare to read from the collection.
	col := reflect.ValueOf(collection)
	// Dereference pointer, if any.
	for col.Kind() == reflect.Ptr {
		col = col.Elem()
	}
	// Make sure at least one record was found.
	if col.Len() == 0 {
		return err404
	}
	// Grab the element. Assume the found value is addressable.
	res := col.Index(0).Elem()

	// Prepare to write to the destination.
	dst := reflect.ValueOf(destination)
	// Make sure the destination is settable.
	if dst.Kind() != reflect.Ptr || !dst.Elem().CanSet() {
		return errorStr("the destination record must be a settable pointer")
	}
	// Set the result.
	dst.Elem().Set(res)

	return nil
}

/****************************** Method Adapters ******************************/

// Returns a datastore key for the given record.
func Key(req *http.Request, record Record) *datastore.Key {
	gc := appengine.NewContext(req)
	return datastore.NewKey(gc, record.Kind(), record.GetId(), 0, nil)
}

// Reads the given record from the Datastore.
func Read(req *http.Request, record Record) error {
	// Check for read permission.
	if !record.Can(req, CodeRead) {
		return err403
	}

	// Read from the Datastore.
	gc := appengine.NewContext(req)
	err := datastore.Get(gc, Key(req, record), record)

	// Compute properties.
	Compute(record)

	// If the record is not found, return a 404 error, otherwise nil.
	if err != nil {
		return err404
	}
	return nil
}

// Saves the given record to the Datastore.
func Save(req *http.Request, record Record) error {
	// If the record is new, check the `create` permission.
	if record.GetId() == "" && !record.Can(req, CodeCreate) {
		return err403
		// Otherwise check for update permission.
	} else if !record.Can(req, CodeUpdate) {
		return err403
	}

	// Validate before saving.
	if len(record.Validate()) != 0 {
		return err400
	}

	// If the id is missing, set a random id.
	if record.GetId() == "" {
		record.SetId(rndId())
	}

	// Save to the Datastore.
	gc := appengine.NewContext(req)
	_, err := datastore.Put(gc, Key(req, record), record)
	return err
}

// Deletes the given record from the Datastore.
func Delete(req *http.Request, record Record) error {
	// Check for delete permission.
	if !record.Can(req, CodeDelete) {
		return err403
	}

	// Delete from the Datastore.
	gc := appengine.NewContext(req)
	err := datastore.Delete(gc, Key(req, record))

	// If deletion fails, assume the record didn't exist and return 404.
	if err != nil {
		return err404
	}
	return nil
}
