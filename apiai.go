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

// Package apiai provides an easy way to handle webhooks coming from api.ai,
// as described in the documentation: https://api.ai/docs/fulfillment
package apiai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Handler provides an easy way to route and handle requests by intent.
type Handler struct {
	mu       sync.RWMutex
	handlers map[string]IntentHandler
}

// NewHandler returns a new empty handler.
func NewHandler() *Handler {
	return &Handler{handlers: make(map[string]IntentHandler)}
}

// Register registers the handler for a given intent.
func (h *Handler) Register(intent string, handler IntentHandler) {
	h.mu.Lock()
	h.handlers[intent] = handler
	h.mu.Unlock()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "could not decode request: "+err.Error(), http.StatusBadRequest)
		return
	}

	intent := req.Result.Metadata.IntentName
	h.mu.RLock()
	handler, ok := h.handlers[intent]
	h.mu.RUnlock()
	if !ok {
		http.Error(w, "could not find handler for "+intent, http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), httpRequestKey, r)
	res, err := handler(ctx, &req)
	if err != nil {
		http.Error(w, "error processing intent:"+err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		http.Error(w, "could not encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", b)
}

// IntentHandler handles an intent.
type IntentHandler func(ctx context.Context, req *Request) (*Response, error)

// A Request contains all of the information to an intent invocation.
type Request struct {
	Lang   string `json:"lang"`
	Status struct {
		ErrorType string `json:"errorType"`
		Code      int    `json:"code"`
	} `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"sessionId"`
	Result    struct {
		Parameters    map[string]string `json:"parameters"`
		Contexts      []struct{}        `json:"contexts"` // TODO
		ResolvedQuery string            `json:"resolvedQuery"`
		Source        string            `json:"source"`
		Score         float64           `json:"score"`
		Speech        string            `json:"speech"`
		Fulfillment   struct {
			Messages []struct {
				// TODO
			} `json:"messages"`
			Speech string `json:"speech"`
		} `json:"fulfillment"`
		ActionIncomplete bool   `json:"actionIncomplete"`
		Action           string `json:"action"`
		Metadata         struct {
			IntentID                  string `json:"intentId"`
			WebhookForSlotFillingUsed bool   `json:"webhookForSlotFillingUsed,string"`
			IntentName                string `json:"intentName"`
			WebhookUsed               bool   `json:"webhookUsed,string"`
		} `json:"metadata"`
	} `json:"result"`
	ID              string `json:"id"`
	OriginalRequest struct {
		Source string                 `json:"source"`
		Data   map[string]interface{} `json:"data"` // TODO
	} `json:"originalRequest"`
}

// Param returns the value associated to the given parameter name.
func (req *Request) Param(name string) string { return req.Result.Parameters[name] }

// A Response is what an intent responds after an invokation.
type Response struct {
	Speech      string `json:"speech"`
	DisplayText string `json:"displayText"`
}

type key string

const httpRequestKey key = "httprequest"

// HTTPRequest returns the HTTP request associated to the given context or nil.
func HTTPRequest(ctx context.Context) *http.Request {
	req, ok := ctx.Value(httpRequestKey).(*http.Request)
	if !ok {
		return nil
	}
	return req
}
