## Description

Database adapter for Golang web applications using the GAE Datastore.

## Features

* Generic functions for collection operations
* Generic methods for record types
* Generic methods for type conversion (records to collections and vice versa)
* Mapping of resource strings to types, resource factory
* Record lifecycle with validation and permission checks
* DB populate helpers

## Contents

* [Description](#description)
* [Features](#features)
* [Installation](#installation)
* [Terminology](#terminology)
* [API Reference](#api-reference)
  * [Record Type](#record-type)
    * [Lifecycle Methods Example](#lifecycle-methods-example)
    * [CRUD Methods Example](#crud-methods-example)
    * [Utility Methods Example](#utility-methods-example)
    * [Key](#key)
    * [Save](#save)
    * [Read](#read)
    * [Delete](#delete)
    * [FindOne](#findone)
  * [Collection Operations](#collection-operations)
    * [Find](#find)
    * [FindAll](#findall)
    * [FindByQuery](#findbyquery)
  * [Permissions](#permissions)
    * [Operation Codes](#operation-codes)
    * [CodeCreate](#codecreate)
    * [CodeRead](#coderead)
    * [CodeUpdate](#codeupdate)
    * [CodeDelete](#codedelete)
  * [Resources](#resources)
    * [Resources](#resources)
    * [NewRecordByResource](#newrecordbyresource)
    * [NewCollectionByResource](#newcollectionbyresource)
    * [SliceOf](#sliceof)
    * [NewRecordFromCollection](#newrecordfromcollection)
  * [Populate](#populate)
    * [RegisterForPopulate](#registerforpopulate)
    * [Populate](#populate)
  * [Setup](#setup)
    * [Config type](#config-type)
    * [Setup](#setup)
  * [Utilities](#utilities)
    * [Compute](#compute)
    * [ErrorCode](#errorcode)
    * [RndId](#rndid)
    * [Log](#log)
    * [ToRecords](#torecords)
  * [Errors](#errors)

## Installation

```shell
go get github.com/Mitranim/gotools/dsadapter
```

In your Go files:

```golang
import (
  dsa "github.com/Mitranim/gotools/dsadapter"
)

var x = dsa.X
```

Optionally, after importing, call `Setup()` to set the options:

```golang
dsadapter.Setup(dsadapter.Config{
  RndId:  dsadapter.RndId, // or your custom id func
  Logger: dsadapter.Log,   // or your custom logging func
})
```

## Terminology

This documentation talks about _records_ and _collections_. A record is a value of a struct type that implements the Record interface. A collection is a slice of that type.

```golang
// record
type Engine struct { /* <...> */ }
engine := &Engine{}

// collection
engines := []*Engine{}
```

## API Reference

This reference is grouped into categories. In addition to listing public types and functions, it includes conventions; in other words, implementation details suggested by the library that need to be met by the user. Read below for details.

### Record type

The generic interface type for a database model.

```golang
type Record interface {
  /* Lifecycle */

  // Validates fields and maps them to error messages.
  Validate() map[string]string
  // Recalculates properties.
  Compute()
  // Permits given operation.
  Can(*http.Request, int) bool

  /* CRUD */

  // Reads self from Datastore by id.
  Read(*http.Request) error
  // Saves self to Datastore.
  Save(*http.Request) error
  // Deletes self from Datastore by id.
  Delete(*http.Request) error

  /* Utilities */

  // Returns own id.
  GetId() string
  // Sets own id to given string.
  SetId(string)
  // Returns own Datastore kind.
  Kind() string
}
```

The library's internal methods depend on the lifecycle and utility methods in the Record interface. They must be implemented by the user in their types; see the examples below.

#### Lifecycle Methods Example

These methods must be implemented by the user. They're automatically called at various points in a record's lifecycle.

##### `Validate() map[string]string`

Must validate own fields and return a map of fields to error messages. `len(errs) == 0` means no error. This method is called by `Record#Save()` before saving to the Datastore.

```golang
func (this *Subscriber) Validate() map[string]string {
  if this.Email == "" {
    return map[string]string{"Email": "Please provide a valid email."}
  }
  return nil
}
```

##### `Compute()`

Provides the opportunity to recalculate computed properties. This method is called automatically when a record is read from the Datastore. When reading a collection, this is called automatically for each record in it.

```golang
func (this *Subscriber) Compute() {
  // Lowercase the email.
  this.Email = strings.ToLower(this.Email)
}
```

##### `Can(*http.Request, int) bool`

Verifies that the user associated with the given request has the rights to perform the given operation. See the [permissions reference](#permissions).

```golang
func (this *Subscriber) Can(req *http.Request, code int) bool {
  // Only allowed to create (gross oversimplification)
  if code == dsadapter.CodeCreate {
    return true
  }
  return false
}
```

#### CRUD Methods Example

`dsadapter` provides generic record methods for three CRUD operations: `Read`, `Save`, and `Delete` (see the reference below). Your types should implement the corresponding methods like so:

```golang
type Engine struct {
  Id   string
  Name string
  // <...>
}

func (this *Engine) Save(req *http.Request) error   { return dsadapter.Save(req, this) }
func (this *Engine) Read(req *http.Request) error   { return dsadapter.Read(req, this) }
func (this *Engine) Delete(req *http.Request) error { return dsadapter.Delete(req, this) }
```

#### Utility Methods Example

These methods must be implemented by your types like so:

```golang
func (this *Engine) Kind() string    { return "Engine" }
func (this *Engine) GetId() string   { return this.Id }
func (this *Engine) SetId(id string) { this.Id = id }
```

They're called in CRUD operations.

#### `Key(\*http.Request, Record) \*datastore.Key`

Creates a Datastore key for the given Record. `Record#Kind()` is called to provide the Datastore kind, and `Record#GetId()` is called to provide the string id for the key. The numeric id associated with a key is always 0. The parent key is always nil.

Each record is identified uniquely in the Datastore by its kind and id within kind.

#### `Save(*http.Request, Record) error`

Generic create/update method for Record types. Saves a record to the Datastore by its kind and id. Example usage:

```golang
func (this *Engine) Save(req *http.Request) error { return dsadapter.Save(req, this) }
```

If the record doesn't have an id (`GetId() == ""`), `dsadapter` generates and assigns a new random id before calling `Key()` and saving the record:

```golang
engine := &Engine{Name: "Zugelgeheiner"}
err := engine.Save(req)

// engine.GetId() -> "3720274029858504238"
```

Be aware that you can't patch a Datastore entity by saving a struct with only _some_ of its fields under the same key. When a struct is created, omitted fields are initialised to zero values. If saved under the same key as an existing entity, it will overwrite it, deleting the existing fields. When updating an entity, you must first read it from the Datastore, update its fields, then save it.

#### `Read(\*http.Request, Record) error`

Generic read method for Record types. Reads a record from the Datastore by its kind and id. Example usage:

```golang
func (this *Engine) Read(req *http.Request) error { return dsadapter.Read(req, this) }

engine := &Engine{Id: "3720274029858504238"}
err := engine.Read(req)

// engine -> {Id: "3720274029858504238", Name: "Zugelgeheiner"}
```

#### `Delete(\*http.Request, Record) error`

Generic delete method for Record types. Deletes a record from the Datastore by its kind and id. Example usage:

```golang
func (this *Engine) Delete(req *http.Request) error { return dsadapter.Delete(req, this) }

engine := &Engine{Id: "3720274029858504238"}
err := engine.Delete(req)
```

#### `FindOne(\*http.Request, Record, map[string]string) error`

Attempts to find one record of the given type by the given parameters and write it to the destination record passed in the function call. The passed record must be a pointer. This is essentially a convenience alias for `FindAll` that writes the result to a record instead of a collection.

Example:

```golang
engine := new(Engine)

err := FindOne(req, engine, map[string]string{"Name": "Zugelgeheiner"})

// engine -> {Id: "3720274029858504238", Name: "Zugelgeheiner"}
```

### Collection Operations

#### `Find(\*http.Request, interface{}, map[string]string, int) error`

Parameters:

```golang
Find(req *http.Request, collection interface{}, params map[string]string, limit int) error
```

Takes a pointer to a collection, a map of filter parameters, and the limit count. Reads the records from the Datastore filtered by these parameters and limited to the given count, writing the result to the collection. Zero or negative limit means no limit.

Example:

```golang
engines := new([]*Engine)

// Suppose we have 10 engines in the Datastore
err := Find(req, engines, nil, 2)

// engines -> &[]*Engine{(*Engine)(0xc2103fa500), (*Engine)(0xc2103fa5a0)}
```

#### `FindAll(\*http.Request, interface{}, map[string]string) error`

Alias of `Find` with 0 limit.

#### `FindByQuery(\*http.Request, interface{})`

Alias of `Find` with 0 limit, where `params` are automatically taken from the `req.URL.Query`.

### Permissions

`dsadapter` checks permissions on each Datastore operation by calling the `Record#Can()` method, passing the http request and the operation code. The implementation of the `Can()` method is up to the user. Generally, the application should check if the user associated with the request has the rights to perform the given operation, possibly depending on the record's relation with other entities, ownership, etc. If the method returns `false`, the CRUD operation is denied and returns an error with the code `403`.

Excessively simplified example:

```golang
// Forbid all but Read operations
func (this *Engine) Can(req *http.Request, code int) bool {
  if code == dsadapter.CodeRead {
    return true
  }
  return false
}

engine := &Engine{Id: "3720274029858504238"}
err := engine.Read(req)
// err == nil  ->  true

err = engine.Delete(req)
// err == nil  ->  false
```

#### Operation Codes

`dsadapter` has four operation codes.

```golang
const (
  CodeCreate = iota
  CodeRead
  CodeUpdate
  CodeDelete
)
```

The codes are [untyped](https://golang.org/ref/spec#Constants) numeric constants.

#### `CodeCreate`

Passed by `Record#Save()` if the record is new (`GetId() == ""`).

#### `CodeRead`

Passed by collection `FindX` functions and by `Record#Read()`.

#### `CodeUpdate`

Passed by `Record#Save()` if the record is not new (`GetId() != ""`).

#### `CodeDelete`

Passed by `Record#Delete()`.

### Resources

`dsadapter` helps you glue together database types and resource URLs.

Resource functions are all about generating new records and collections in a generic way. They let you allocate new records and collections by the given resource string, get collection by record, and vice versa. All of them are generic, that is, they don't force you to deal with concrete types.

Types are not first-class values in golang. We emulate first-class types with nil pointers to values of those types. A nil pointer carries only its type information, is immutable, and can be used as a reference to make a new value of its type. The syntax looks like this:

```golang
// Concrete type
type Quasar struct {
  // <...>
}
// Nil pointer
(*Quasar)(nil)
```

#### `Resources map[string]Record`

Package-wide map of resource strings to record types, where types are represented with nil pointers to values. It's used internally by other resource methods. Register your resources by assigning them to this map. Example:

```golang
Resources["engines"] = (*Engine)(nil)
```

#### `NewRecordByResource(string) Record`

Returns a pointer to a newly allocated zero value of the concrete type associated with the resource.

The gotcha here is that while the underlying value has a concrete type, the returned value has the interface type `Record`. To get the underlying value, use a type assertion.

Example:

```golang
Resources["engines"] = (*Engine)(nil)

// Get a new engine as Record
engine := NewRecordByResource("engines")

// The underlying value is a concrete *Engine
// fmt.Printf("#%v", engine)  ->  &Engine{Id: "", Name: ""}

// Get a concrete Engine
eng := engine.(*Engine)

// Check for another type with a type assertion
_, ok := engine.(*Quasar)
// ok -> false
```

#### `NewCollectionByResource(string) interface{}`

Similar to `NewRecordByResource`, but instead of creating a record of the concrete resource type, it creates an empty slice of that type. Returns a pointer to the slice.

The gotcha here is that the returned pointer has type `interface{}`. [This article](https://github.com/golang/go/wiki/InterfaceSlice) explains why. It can still be passed into functions that accept collections as `interface{}`, like `FindByQuery`. If you need a slice of records, call `ToRecords()` to convert it to the `[]Record` type. If you want the concrete underlying type, do a type assertion.

Example:

```golang
Resources["engines"] = (*Engine)(nil)

// Get a new collection as interface{}
engines := NewCollectionByResource("engines")

// The underlying value is a concrete *[]*Engine
// fmt.Printf("#%v", engine)  ->  &[]*Engine{}

// Get the concrete engines with a type assertion
engs, ok := engines.(*[]*Engine)
```

#### `SliceOf(interface{}) interface{}`

Allocates an empty non-nil slice of the given value's concrete type and returns a pointer to it masquerading as `interface{}`. This is used internally in `NewCollectionByResource`.

#### `NewRecordFromCollection(interface{}) (Record, error)`

Takes a pointer to a collection, allocates a new empty Record of the same concrete type, and returns a pointer to it. This is used internally in collection `FindX` functions to get a record to query its `Kind()` and `Can()` methods.

### Populate

`dsadapter` comes with a primitive populate routine to help populate the database with data mockups.

#### `RegisterForPopulate(interface{})`

Takes a collection, converts it to the `[]Record` type, and registers a function that repopulates this Datastore kind with the given collection when `Populate()` is called. Any existing records of this kind are deleted before saving new ones. If an error occurs at any stage, this particular repopulate func is aborted.

Example:

```golang
engines := []*Engine{engine0, engine1}

RegisterForPopulate(engines)
```

Only one populate per kind can be registered; attempts to register more than one collection of the same kind are silently ignored.

#### `Populate(*http.Request)`

Calls each populate function registered with RegisterForPopulate.

### Setup

After importing `dsadapter`, you can optionally reconfigure it by calling `Setup()` and passing a configuration struct Config with the appropriate options.

#### Config type

```golang
type Config struct {
  // Function to call when generating a missing id for a new record. If omitted,
  // the default dsadapter.RndId function is used. To disable automatic id generation
  // (not recommended), pass a function that returns an empty string.
  RndId func() string

  // Logger function to call on populate and critical errors. If omitted, no
  // logging is done. Pass dsadapter.Log to use the default logger.
  Logger func(...interface{})
}
```

#### `Setup(Config) error`

Saves the passed config into a package-wide global. The options take effect immediately.

### Utilities

#### `Compute(interface{})`

Takes any value and calls its `Compute()` method, if available. If the value is a slice, loops over it and calls the `Compute()` method on each element, if possible. This is called automatically after reading a record or collection from the Datastore to recalculate computed properties.

#### `ErrorCode(error) int`

Reads the error message and returns the number from its beginning, if it falls in the http error code range: `400 <= x <= 599`. If the error doesn't begin with a number or it falls outside the boundary, the function returns `500`. If the error is nil, the function returns `200`.

Example:

```golang
// From utils-private.go
const err404 = errorStr("404 not found")
ErrorCode(err404) // -> 404

var unknownErr = errors.New("no squirrel-to-squid interface defined")
ErrorCode(unknownErr) // -> 500
```

When using functions and methods backed by `dsadapter` in your request handlers, you should examine them with `ErrorCode()` and set the appropriate http status code in your response.

#### `RndId() string`

Generates a random string id. This is used by default to make random ids for new records when saving them. You can override it by passing a custom `RndId` value in a `Setup()` call.

#### `Log(...interface{})`

A simple logging function. Alias for `println` that automatically does `"%v"` on each value. Pass it into a `Setup()` call to use the default logger (logging is disabled by default).

#### `ToRecords(interface{}) []Record`

Takes a slice of any type and converts it to a slice of records. If the value is not a slice, this returns nil. Non-Record values are ignored. If the original slice didn't contain any Records, the result will be zero length. This is used internally in `RegisterForPopulate()` to convert the given collection to a slice of records.

Example:

```golang
// Concrete slice of Engines
engines := []*Engine{engine0, engine1}

// Interface version
engi := interface{}(engines)

// []Record version
engs := ToRecords(engi)
```

### Errors

Most errors generated by `dsadapter` begin with an http status code corresponding to the error type. When handling a `dsadapter` error, you should examine it with `ErrorCode()` and set the resulting status code in your http response writer.

Example:

```golang
engine := &Engine{Id: "nonexistend id"}
err := engine.Read(req)

// err.Error() -> "404 not found"
// dsadapter.ErrorCode(err) -> 404
```

Quick reference:
* failed `Validate()` → 400
* failed `Can()` → 403
* failed `Read()` or `Delete()` → 404

Some errors generated by the Datastore are returned as-is. `ErrorCode()` returns `500` for them.
