package dsadapter

// Database populate utilities.

import (
	// Standard
	"net/http"
)

// Returns the map of Datastore kind strings to populate funcs registered with
// RegisterForPopulate. Creates the map if it's nil.
func (this *StateInstance) PopulateFuncs() map[string]func(*http.Request) {
	if this.populateFuncs == nil {
		this.populateFuncs = map[string]func(*http.Request){}
	}
	return this.populateFuncs
}

// Loops over each registered populate func and calls it. The App Engine library
// panics if we do this asynchronously with goroutines, thus the synchrony. We
// could probably get around this by passing around the same GAE context.
func (this *StateInstance) Populate(req *http.Request) {
	for _, fn := range this.PopulateFuncs() {
		fn(req)
	}
}

// Registers the given records for populate.
func (this *StateInstance) RegisterForPopulate(values interface{}) {
	// Convert to the []Record type.
	records := ToRecords(values)

	// Make sure we have at least one record.
	if len(records) == 0 {
		return
	}

	// Ignore if there's already a function under this collection's kind.
	kind := records[0].Kind()
	if this.PopulateFuncs()[kind] != nil {
		return
	}

	// Register a populate func.
	this.PopulateFuncs()[kind] = func(req *http.Request) {
		this.log(req, "   populating kind:", kind)

		// Retrieve all existing records of this kind to delete them.
		oldRecs := this.SliceOf(records[0])
		err := this.FindAll(req, oldRecs, nil)
		if err != nil {
			this.log(req, "!! unexpected error when retrieving old records during populate:", err)
		}

		// Loop over and call the Delete method of each old record.
		func() {
			for _, rec := range ToRecords(oldRecs) {
				// Try to delete; abort the sequence if this fails.
				err := rec.Delete(req)
				if err != nil {
					this.log(req, "!! unexpected error when trying to delete an old record during populate:", err)
					return
				}
			}
			this.log(req, "-- deleted all records of kind:", kind)
		}()

		// Loop over records and save them.
		func() {
			for _, record := range records {
				// This point in a record's lifecycle is equivalent to it being created
				// from json or read from a database, so we must call its Compute method.
				this.Compute(record)
				// Try to save it and abort the sequence if this fails.
				if err := record.Save(req); err != nil {
					this.log(req, "!! failed to save record during populate:", err)
					this.log(req, "!! aborting populate of kind:", kind)
					return
				}
			}
			this.log(req, "++ successfully populated all records of kind:", kind)
		}()
	}
}
