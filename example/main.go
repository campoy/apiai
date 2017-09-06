// Copyright 2017 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// This example shows how to create a simple chat bot that reads a number from
// the requests and doubles it.
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
	num, err := strconv.Atoi(req.Param("num"))
	if err != nil {
		return nil, fmt.Errorf("could not parse number %q: %v", req.Param("number"), err)
	}
	return &apiai.Response{
		Speech: fmt.Sprintf("%d times two equals %d", num, 2*num),
	}, nil
}
