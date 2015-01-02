package context

// Example middleware snippet: how to map context to a Martini request handler.

/*

import (
  // Standard
  "net/http"

  // Third party
  context "github.com/Mitranim/gotools/context"
  "github.com/go-martini/martini"
)

// Maps a new context to a martini request handler chain.
func MapContext(rw http.ResponseWriter, req *http.Request, cont martini.Context) {
  // Context object.
  ct := context.NewContext(rw, req, context.Config{Render: MyRenderFunc})
  // Mapping.
  cont.MapTo(ct, (*context.Context)(nil))
}

*/
