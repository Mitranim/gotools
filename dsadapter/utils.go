package dsadapter

// Public utilities and constants.

import (
	// Standard
	"math/rand"
	"net/http"
	"strconv"
	"time"

	// App Engine
	"appengine"

	// Third party
	"github.com/Mitranim/gotools/utils"
)

// Weak-entropy seed.
func init() {
	rand.Seed(time.Now().UnixNano())
}

/********************************* Constants *********************************/

// CRUD operation codes.
const (
	CodeCreate = iota
	CodeRead
	CodeUpdate
	CodeDelete
)

/********************************* Utilities *********************************/

// Default id function. Makes a random unique id.
func RndId() string {
	i := rand.Int()
	return strconv.Itoa(i)
}

// Republish the error-to-code converter.
var ErrorCode = utils.ErrorCode

// Logs to a GAE context with level `debug`.
func Log(req *http.Request, values ...interface{}) {
	gc := appengine.NewContext(req)
	gc.Debugf(repeat("%v", len(values)), values...)
}
