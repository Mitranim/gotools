package dsadapter

// Database populate utilities.

import (
	// Standard
	"net/http"
)

/********************************** Globals **********************************/

var PopulateFuncs = map[string]func(*http.Request){}

/********************************* Utilities *********************************/

// Loops over each registered populate func and calls it. The App Engine library
// panics if we do this asynchronously with goroutines, thus the synchrony. We
// could probably get around this by passing around the same GAE context.
func Populate(req *http.Request) {
	for _, fn := range PopulateFuncs {
		fn(req)
	}
}

// Registers the given records for populate.
func RegisterForPopulate(values interface{}) {
	// Convert to the []Record type.
	records := ToRecords(values)

	// Make sure we have at least one record.
	if len(records) == 0 {
		return
	}

	// Ignore if there's already a function under this collection's kind.
	kind := records[0].Kind()
	if PopulateFuncs[kind] != nil {
		return
	}

	// Register a populate func.
	PopulateFuncs[kind] = func(req *http.Request) {
		log(req, "   populating kind:", kind)

		// Retrieve all existing records of this kind to delete them.
		oldRecs := SliceOf(records[0])
		err := FindAll(req, oldRecs, nil)
		if err != nil {
			log(req, "!! unexpected error when retrieving old records during populate:", err)
		}

		// Loop over and call the Delete method of each old record.
		func() {
			for _, rec := range ToRecords(oldRecs) {
				// Try to delete; abort the sequence if this fails.
				err := rec.Delete(req)
				if err != nil {
					log(req, "!! unexpected error when trying to delete an old record during populate:", err)
					return
				}
			}
			log(req, "-- deleted all records of kind:", kind)
		}()

		// Loop over records and save them.
		func() {
			for _, record := range records {
				// This point in a record's lifecycle is equivalent to it being created
				// from json or read from a database, so we must call its Compute method.
				Compute(record)
				// Try to save it and abort the sequence if this fails.
				if err := record.Save(req); err != nil {
					log(req, "!! failed to save record during populate:", err)
					log(req, "!! aborting populate of kind:", kind)
					return
				}
			}
			log(req, "++ successfully populated all records of kind:", kind)
		}()
	}
}
