[![GoDoc](https://godoc.org/github.com/campoy/apiai?status.svg)](http://godoc.org/github.com/campoy/apiai)
[![Build Status](https://travis-ci.org/campoy/apiai.svg)](https://travis-ci.org/campoy/apiai) [![Go Report Card](https://goreportcard.com/badge/github.com/campoy/apiai)](https://goreportcard.com/report/github.com/campoy/apiai)

# apiai

This package provides an easy way to handle webhook requests for
[api.ai](https://api.ai) fulfillments.

You can find more in the [fulfillment docs](https://api.ai/docs/fulfillment).

This is totally experimental for now, so please file let me know what you think about it.
File issues!

## example

This is a simple example fo an application handling a single intent called number,
which it simply doubles as a result.

[embedmd]:# (example/main.go /package main/ $)
```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/campoy/apiai"
)

func main() {
	h := apiai.NewHandler()
	h.Register("double", doubleHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", h))
}

func doubleHandler(ctx context.Context, req *apiai.Request) (*apiai.Response, error) {
	num, err := strconv.Atoi(req.Param("number"))
	if err != nil {
		return nil, fmt.Errorf("could not parse number %q: %v", req.Param("number"), err)
	}
	return &apiai.Response{
		Speech: fmt.Sprintf("%d times two equals %d", num, 2*num),
	}, nil
}
```

### Disclaimer

This is not an official Google product (experimental or otherwise), it is just
code that happens to be owned by Google.
