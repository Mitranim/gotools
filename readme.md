### Description

Micro-framework for Golang web servers. Condenses many common request handler tasks to a single function call or line of code.

Has modules for smart page rendering, contextual handler tasks, and database modeling.

Although the `gotools` are tied together in the root package, each component is independent from others and can be used in isolation. See their respective docs:
* [`render` readme](render)
* [`context` readme](context)
* [`dsadapter` readme](dsadapter)

`gotools` are orthogonal to middleware frameworks like [Martini](https://github.com/go-martini/martini) or [Gorilla](http://www.gorillatoolkit.org), and should be combined with them. The context component includes an [example code snippet](context/middleware.go) how to plug it into Martini as middleware.

### Installation

```shell
go get github.com/Mitranim/gotools
```

In your Go files:

```golang
import (
  gt "github.com/Mitranim/gotools"
)
```

The root `gotools` package imports all of its components and republishes their parts with some adaptations. When using the whole package, you only need to import `gotools`.
