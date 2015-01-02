package gotools

// Imports and republishes the components, wrapping them into one package.

import (
	// Standard
	"net/http"

	// Components
	"github.com/Mitranim/gotools/context"
	"github.com/Mitranim/gotools/dsadapter"
	"github.com/Mitranim/gotools/render"

	// Utilities
	"github.com/Mitranim/gotools/utils"
)

/********************************** Shared ***********************************/

// Functions
var (
	ErrorCode = utils.ErrorCode
	CodePath  = utils.CodePath
	ErrorPath = utils.ErrorPath
	Log       = utils.Log
)

/********************************** context **********************************/

// Functions
var (
	Panic   = context.Panic
	Recover = context.Recover
)

// Types
type Context context.Context
type ContextConfig context.Config

// Adapters
func NewContext(rw http.ResponseWriter, req *http.Request, config ContextConfig) Context {
	return context.NewContext(rw, req, context.Config(config))
}

/********************************** render ***********************************/

// Types
type RenderState render.State
type RenderConfig render.Config

// Adapters
func RenderSetup(config RenderConfig) (RenderState, error) {
	state, err := render.Setup(render.Config(config))
	return RenderState(state), err
}

/********************************* dsadapter *********************************/

// Functions
var (
	DsaLog = dsadapter.Log
	RndId  = dsadapter.RndId
)

// Constants
const (
	CodeCreate = dsadapter.CodeCreate
	CodeRead   = dsadapter.CodeRead
	CodeUpdate = dsadapter.CodeUpdate
	CodeDelete = dsadapter.CodeDelete
)

// Types
type Record dsadapter.Record
type DsaConfig dsadapter.Config
type DsaState dsadapter.State

// Adapters
func DsaSetup(config DsaConfig) DsaState {
	return dsadapter.Setup(dsadapter.Config(config))
}
