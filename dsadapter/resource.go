package dsadapter

// Utilities for mapping types to URL resources and generating derived objects.

import (
	// Standard
	"reflect"

	// Third party
	"github.com/Mitranim/gotools/utils"
)

/**
 * Returns the map of URL resource strings to record types, creating it if
 * it's nil. Typical usage:
 *
 *   // Define a type.
 *   type User struct { <...> }
 *
 *   // Register the type as a resource.
 *   dsadapter.Resources()["users"] = (*User)(nil)
 */
func (this *stateInstance) Resources() map[string]Record {
	if this.resources == nil {
		this.resources = map[string]Record{}
	}
	return this.resources
}

// If there is a record type registered under the given resource name, this
// allocates a zero value of that type and returns a pointer to it masquerading
// as Record. If there isn't a matching type, this returns nil.
func (this *stateInstance) NewRecordByResource(name string) Record {
	// Grab a reference record.
	record := this.Resources()[name]
	if record == nil {
		return nil
	}
	// Make a copy.
	val := reflect.New(reflect.TypeOf(record).Elem())
	return val.Interface().(Record)
}

// If there is a record type registered under the given resource name, this
// allocates a nil slice of that type and returns a pointer to it masquerading
// as interface{}. If there isn't a matching type, this returns nil.
func (this *stateInstance) NewCollectionByResource(name string) interface{} {
	// Grab a reference record.
	record := this.Resources()[name]
	if record == nil {
		return nil
	}
	// Make a slice of it and return a pointer to that slice.
	return this.SliceOf(record)
}

// Allocates a zero-length non-nil slice of the given value's type, takes its
// pointer, and returns the pointer masquerading as interface{}.
func (this *stateInstance) SliceOf(value interface{}) interface{} {
	return reflect.New(reflect.SliceOf(reflect.TypeOf(value))).Interface()
}

// Takes a pointer to a collection and returns a new Record of its type. Use
// SliceOf for the (roughly) opposite effect.
func (this *stateInstance) NewRecordFromCollection(collection interface{}) (Record, error) {
	// We're going to return this error if anything goes wrong.
	err := utils.Error("a collection must be a slice of a struct pointer type that implements Record")

	val := reflect.ValueOf(collection)

	// Make sure this is a pointer and dereference it.
	if val.Kind() != reflect.Ptr {
		return nil, err
	}
	val = val.Elem()

	// Make sure it's a slice and get its element type.
	if val.Kind() != reflect.Slice {
		return nil, err
	}
	elemType := val.Type().Elem()

	// Make sure the element type is a struct type or struct pointer type.
	isStruct := false
	if elemType.Kind() == reflect.Struct {
		isStruct = true
	} else if elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct {
		isStruct = true
	}
	if !isStruct {
		return nil, err
	}

	// Make an example Record.
	rec := reflect.New(elemType).Elem().Interface()

	// Make sure it implements the Record interface.
	record, ok := rec.(Record)
	if !ok {
		return nil, err
	}

	// Return the result.
	return record, nil
}
