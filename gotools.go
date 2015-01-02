package gotools

// Imports and republishes the components, wrapping them into one package.

import (
	// Components
	"github.com/Mitranim/gotools/context"
	"github.com/Mitranim/gotools/dsadapter"
	"github.com/Mitranim/gotools/render"

	// Utilities
	"github.com/Mitranim/gotools/utils"
)

/********************************** Shared ***********************************/

var (
	ErrorCode = utils.ErrorCode
	CodePath  = utils.CodePath
	ErrorPath = utils.ErrorPath
	Log       = utils.Log
)

/********************************** context **********************************/

/*-------------------------------- Functions --------------------------------*/

var (
	NewContext = context.NewContext
	Panic      = context.Panic
	Recover    = context.Recover
)

/*---------------------------------- Types ----------------------------------*/

type Context context.Context
type ConfigContext context.Config
type ContextInstance context.ContextInstance

/*-------------------------------- Adapters ---------------------------------*/

func SetupContext(config ConfigContext) error {
	return context.Setup(context.Config(config))
}

/********************************** render ***********************************/

/*-------------------------------- Functions --------------------------------*/

var (
	Render           = render.Render
	RenderError      = render.RenderError
	RenderPage       = render.RenderPage
	RenderStandalone = render.RenderStandalone
	InlineFiles      = render.InlineFiles
	Pages            = render.Pages
	Standalone       = render.Standalone
)

/*---------------------------------- Types ----------------------------------*/

type ConfigRender render.Config

/*-------------------------------- Adapters ---------------------------------*/

func SetupRender(config ConfigRender) error {
	return render.Setup(render.Config(config))
}

/********************************* dsadapter *********************************/

/*-------------------------------- Functions --------------------------------*/

var (
	Compute                 = dsadapter.Compute
	Delete                  = dsadapter.Delete
	Debug                   = dsadapter.Debug
	Find                    = dsadapter.Find
	FindAll                 = dsadapter.FindAll
	FindByQuery             = dsadapter.FindByQuery
	FindOne                 = dsadapter.FindOne
	Key                     = dsadapter.Key
	NewCollectionByResource = dsadapter.NewCollectionByResource
	NewRecordByResource     = dsadapter.NewRecordByResource
	NewRecordFromCollection = dsadapter.NewRecordFromCollection
	Populate                = dsadapter.Populate
	PopulateFuncs           = dsadapter.PopulateFuncs
	Read                    = dsadapter.Read
	RegisterForPopulate     = dsadapter.RegisterForPopulate
	RndId                   = dsadapter.RndId
	Save                    = dsadapter.Save
	SliceOf                 = dsadapter.SliceOf
	ToRecords               = dsadapter.ToRecords
)

/*--------------------------------- Values ----------------------------------*/

const (
	CodeCreate = dsadapter.CodeCreate
	CodeRead   = dsadapter.CodeRead
	CodeUpdate = dsadapter.CodeUpdate
	CodeDelete = dsadapter.CodeDelete
)

var Resources = dsadapter.Resources

/*---------------------------------- Types ----------------------------------*/

type Record dsadapter.Record
type ConfigDsa dsadapter.Config

/*-------------------------------- Adapters ---------------------------------*/

func SetupDsa(config ConfigDsa) error {
	return dsadapter.Setup(dsadapter.Config(config))
}
