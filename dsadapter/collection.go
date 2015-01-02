package dsadapter

// Collection related functions.

import (
	// Standard
	"net/http"

	// App Engine
	"appengine"
	"appengine/datastore"
)

/****************************** Query Functions ******************************/

// Takes a pointer to a Collection and finds records for it, filtered by the
// given params and limited to the given count. The records are added to the
// collection. The collection may be created with reflect like so:
// reflect.New(<slice type>).Interface(). Zero or negative limit means no limit.
func Find(req *http.Request, collection interface{}, params map[string]string, limit int) error {
	gc := appengine.NewContext(req)

	// Make a Record of this collection's type to get its Datastore kind.
	record, err := NewRecordFromCollection(collection)
	if err != nil {
		return err
	}

	// Check for read permission.
	if !record.Can(req, CodeRead) {
		return err403
	}

	// Form a query.
	q := datastore.NewQuery(record.Kind())

	// Apply params, if any.
	for key, param := range params {
		q = q.Filter(key+" =", param)
	}

	// Apply limit, if any. Zero or negative means no limit.
	if limit > 0 {
		q = q.Limit(limit)
	}

	// Run the query, writing to the collection.
	_, err = q.GetAll(gc, collection)

	if err != nil {
		log("-- error in datastore query:", err)
		return err
	}

	// Compute properties on children.
	Compute(collection)

	return nil
}

// Takes a pointer to a Collection and finds records for it, filtered by the
// given params.
func FindAll(req *http.Request, collection interface{}, params map[string]string) error {
	return Find(req, collection, params, 0)
}

// Takes a pointer to a Collection and finds records for it, filtered by the URL
// query params (if any).
func FindByQuery(req *http.Request, collection interface{}) error {
	return Find(req, collection, toParams(req.URL.Query()), 0)
}
