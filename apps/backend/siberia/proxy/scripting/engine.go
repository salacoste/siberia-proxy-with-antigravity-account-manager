package scripting

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/dop251/goja"
)

// ScriptEngine manages the execution of JavaScript hooks
type ScriptEngine struct {
	vmPool sync.Pool
	Script string // The user-defined script
	Active bool
}

func NewScriptEngine() *ScriptEngine {
	return &ScriptEngine{
		vmPool: sync.Pool{
			New: func() interface{} {
				return goja.New()
			},
		},
		Active: false,
	}
}

// UpdateScript updates the current script
func (e *ScriptEngine) UpdateScript(script string, active bool) {
	e.Script = script
	e.Active = active
}

// JsRequest is the object exposed to JS
type JsRequest struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

// JsResponse is the object exposed to JS
type JsResponse struct {
	StatusCode int                 `json:"status"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
}

// RunOnRequest executes the onRequest function in the script
func (e *ScriptEngine) RunOnRequest(req *http.Request) (*http.Request, error) {
	if !e.Active || e.Script == "" {
		return req, nil
	}

	vm := e.vmPool.Get().(*goja.Runtime)
	defer e.vmPool.Put(vm)

	//Reset VM? Goja doesn't have a clear "Reset".
	//Ideally we use a fresh VM or standard loop.
	//For safety in V1, let's just make a New() one if we face pollution issues,
	//but for simple hooks, sharing via Pool is "okay" if we don't leak global state.
	//Actually, users might define globals. A pool might leak state between requests.
	//Safest is NEW VM per request. Goja is fast enough.
	//Let's discard the pool idea for now to guarantee isolation.
	vm = goja.New()

	// 1. Setup Environment
	_, err := vm.RunString(e.Script)
	if err != nil {
		return req, fmt.Errorf("script syntax error: %w", err)
	}

	onRequest, ok := goja.AssertFunction(vm.Get("onRequest"))
	if !ok {
		// No onRequest hook defined, pass through
		return req, nil
	}

	// 2. Prepare Context
	bodyStr := ""
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(b))
		bodyStr = string(b)
	}

	jsReq := &JsRequest{
		Method:  req.Method,
		URL:     req.URL.String(),
		Headers: req.Header,
		Body:    bodyStr,
	}

	// 3. Execute
	val, err := onRequest(goja.Undefined(), vm.ToValue(jsReq))
	if err != nil {
		return req, fmt.Errorf("runtime error: %w", err)
	}

	// 4. Apply Changes (if object returned)
	// We expect the function to return the modified request object, or we modify the one passed in?
	// JS is pass-by-reference for objects. We can inspect jsReq back?
	// Goja maps struct fields to JS properties. If script modifies properties, does it reflect in struct?
	// No, vm.ToValue creates a copy/proxy.
	// We should return the object from JS.

	// Let's assume user does: return req;
	resObj := val.ToObject(vm)

	// Reconstruct http.Request
	// Method
	if v := resObj.Get("Method"); v != nil {
		req.Method = v.String()
	}
	// URL - complicated to parse back? Assuming mostly Method/Headers/Body changes.
	// Skipping URL rewrite for now in V1 unless crucial.

	// Headers
	if v := resObj.Get("Headers"); v != nil {
		// Convert back to map[string][]string
		// This is tricky with Goja.
		// Simplified: We assume user modified the inputs.
		// Let's Re-Export to struct.
		var newJsReq JsRequest
		err = vm.ExportTo(val, &newJsReq)
		if err == nil {
			req.Method = newJsReq.Method
			req.Header = newJsReq.Headers
			if newJsReq.Body != bodyStr {
				req.Body = io.NopCloser(bytes.NewBufferString(newJsReq.Body))
				req.ContentLength = int64(len(newJsReq.Body))
			}
		}
	}

	return req, nil
}

// RunOnResponse executes the onResponse function in the script
func (e *ScriptEngine) RunOnResponse(resp *http.Response) (*http.Response, error) {
	if !e.Active || e.Script == "" {
		return resp, nil
	}

	vm := goja.New()
	_, err := vm.RunString(e.Script)
	if err != nil {
		return resp, fmt.Errorf("script syntax error: %w", err)
	}

	onResponse, ok := goja.AssertFunction(vm.Get("onResponse"))
	if !ok {
		return resp, nil
	}

	// Read Body
	bodyStr := ""
	if resp.Body != nil {
		b, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(b))
		bodyStr = string(b)
	}

	jsResp := &JsResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       bodyStr,
	}

	val, err := onResponse(goja.Undefined(), vm.ToValue(jsResp))
	if err != nil {
		return resp, fmt.Errorf("runtime error: %w", err)
	}

	// Apply Changes
	var newJsResp JsResponse
	err = vm.ExportTo(val, &newJsResp)
	if err == nil {
		resp.StatusCode = newJsResp.StatusCode
		resp.Status = fmt.Sprintf("%d Script Modified", newJsResp.StatusCode)
		resp.Header = newJsResp.Headers
		if newJsResp.Body != bodyStr {
			resp.Body = io.NopCloser(bytes.NewBufferString(newJsResp.Body))
			resp.ContentLength = int64(len(newJsResp.Body))
		}
	}

	return resp, nil
}
