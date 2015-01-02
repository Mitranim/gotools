### Description

Micro-framework for Golang web servers. Condenses many common request handler tasks to a single function call or line of code.

Has modules for smart page rendering, contextual handler tasks, and database modeling.

Each module is completely independent from others and can be used in isolation. See their respective docs:
* [render](render)
* [context](context)
* [dsadapter](dsadapter)

`gotools` are orthogonal to middleware frameworks like [Martini](https://github.com/go-martini/martini) or [Gorilla](http://www.gorillatoolkit.org), and should be combined with them. The context component includes an [example code snippet](context/middleware.go) how to plug it into Martini as middleware.

### Installation

```shell
go get github.com/Mitranim/gotools
```

In your Go files:

```golang
import (
  "github.com/Mitranim/gotools"
)
```
