// api.go
package wrapper

import (
	"encoding/json"
	"net/http"
	"strings"

	"honnef.co/go/augeas"
)

// This type acts as a REST API server adapter for the augeas.Augeas type.
type ConfigWrapper struct {
	augeas.Augeas
}

// This structure represents a displayable Augeas API error message.
type SerializableError struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	MinorMessage string `json:"minor_message"`
	Details      string `json:"details"`
}

// This type maps the parameter names to their values.
type ParamMap map[string]string

// This type maps the names of the REST API objects to their values.
type ObjectMap map[string]interface{}

// This function represents a handler for a custom request.
type RequestHandler func(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error)

// New creates a new Augeas REST API server.
func New(configRoot string, loadPath string, flags augeas.Flag) (ConfigWrapper, error) {
	aug, err := augeas.New(configRoot, loadPath, flags)
	return ConfigWrapper{aug}, err
}

func getSerializableError(err augeas.Error) *SerializableError {
	return &SerializableError{
		Code:         int(err.Code),
		Message:      err.Message,
		MinorMessage: err.MinorMessage,
		Details:      err.Details,
	}
}

func handleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch err := err.(type) {
	case augeas.Error:
		var statusCode int

		switch err.Code {
		case augeas.CouldNotInitialize:
			statusCode = http.StatusServiceUnavailable
		case augeas.NoMatch:
			statusCode = http.StatusNotFound
		case augeas.ENOMEM:
			statusCode = http.StatusServiceUnavailable
		case augeas.EINTERNAL:
			statusCode = http.StatusInternalServerError
		case augeas.EPATHX:
			statusCode = http.StatusBadRequest
		case augeas.ENOMATCH:
			statusCode = http.StatusNotFound
		case augeas.EMMATCH:
			statusCode = http.StatusRequestedRangeNotSatisfiable
		case augeas.ESYNTAX:
			statusCode = http.StatusInternalServerError
		case augeas.ENOLENS:
			statusCode = http.StatusInternalServerError
		case augeas.EMXFM:
			statusCode = http.StatusTooManyRequests
		case augeas.ENOSPAN:
			statusCode = http.StatusUnsupportedMediaType
		case augeas.EMVDESC:
			statusCode = http.StatusUnprocessableEntity
		case augeas.ECMDRUN:
			statusCode = http.StatusInternalServerError
		case augeas.EBADARG:
			statusCode = http.StatusBadRequest
		}

		w.WriteHeader(statusCode)
		res, _ := json.Marshal(getSerializableError(err))
		w.Write(res)
	}
}

// HandleRequest is used to invoke the request handler for the respective request.
func (wrapper ConfigWrapper) HandleRequest(w http.ResponseWriter, r *http.Request, method string, prefix string, handler RequestHandler) {
	w.Header()["Content-Type"] = []string{"application/json"}

	if r.Method != method {
		handleError(w, &augeas.Error{Code: augeas.EBADARG})
		return
	}

	object, err := handler(wrapper, "/"+strings.TrimPrefix(r.URL.Path, prefix), r)

	if err != nil {
		handleError(w, err)
		return
	}

	res, _ := json.Marshal(object)
	w.Write(res)
}

// ParseRequestBodyParams parses the parameters from the JSON request string.
func ParseRequestBodyParams(r *http.Request) (ParamMap, error) {
	var body []byte
	r.Body.Read(body)
	params := ParamMap{}
	err := json.Unmarshal(body, params)
	return params, err
}

// List retrieves the list of paths to the direct children of the Augeas node
// with the given path as a JSON array.
func (wrapper ConfigWrapper) List(path string) ([]string, error) {
	return wrapper.Match(path + "*")
}

// WebList retrieves the list of paths to the direct children of the Augeas node
// with the given path as a JSON array.
func WebList(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	return wrapper.List(path)
}

// WebMatch retrieves the list of paths to the Augeas nodes matching the given
// pattern string as a JSON array.
// The pattern string's special syntax is as follows: the `*` character represents
// any character in a label, and `[i]` represents the i-th element of a node array.
func WebMatch(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	return wrapper.Match(path)
}

// WebGet retrieves the value of the Augeas node with the given path as a JSON
// object with one property: `value`.
func WebGet(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	value, err := wrapper.Get(path)
	return ObjectMap{"value": value}, err
}

// WebGetAll retrieves the values of all Augeas nodes with the given path (i.e. all
// nodes from the array with that path) as a JSON array.
func WebGetAll(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	return wrapper.GetAll(path)
}

// WebGetLabel retrieves the label (the last component of the path) of the Augeas
// node with the given path.
func WebGetLabel(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	return wrapper.Label(path)
	value, err := wrapper.Label(path)
	return ObjectMap{"label": value}, err
}

// WebSet assigns the value passed as the JSON `value` property in the request
// string to the Augeas node with the given path.
// Its parent nodes are created as necessary if they do not exist.
func WebSet(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	params, err := ParseRequestBodyParams(r)
	if err != nil {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	value, ok := params["value"]
	if !ok {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	err = wrapper.Set(path, value)
	return ObjectMap{}, err
}

// WebSetMultiple assigns the value passed as the JSON `value` property in the
// request string to all nodes in the Augeas node array with the given path whose
// relative paths to the array match the pattern passed as the JSON `pattern`
// property.
// If the pattern is an empty string, all nodes in the array are matched.
func WebSetMultiple(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	params, err := ParseRequestBodyParams(r)
	if err != nil {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	pattern, ok := params["pattern"]
	if !ok {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	value, ok := params["value"]
	if !ok {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	num, err := wrapper.SetMultiple(path, pattern, value)
	return ObjectMap{"count": num}, err
}

// WebClear clears (i.e. sets to NULL) the value of the Augeas node with the given
// path.
// Its parent nodes are created as necessary if they do not exist.
func WebClear(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	err := wrapper.Clear(path)
	return ObjectMap{}, err
}

// WebInsertBefore inserts a new Augeas node with a label passed as the JSON
// `label` property in the request string before the node with the given path.
func WebInsertBefore(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	params, err := ParseRequestBodyParams(r)
	if err != nil {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	label, ok := params["label"]
	if !ok {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	err = wrapper.Insert(path, label, true)
	return ObjectMap{}, err
}

// WebInsertAfter inserts a new Augeas node with a label passed as the JSON
// `label` property in the request string after the node with the given path.
func WebInsertAfter(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	params, err := ParseRequestBodyParams(r)
	if err != nil {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	label, ok := params["label"]
	if !ok {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	err = wrapper.Insert(path, label, false)
	return ObjectMap{}, err
}

// WebRemove removes the whole Augeas subtree with the given path, including all of
// its descendants.
func WebRemove(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	num := wrapper.Remove(path)
	return ObjectMap{"count": num}, nil
}

// WebMove moves the whole Augeas subtree with the path passed as the JSON `source`
// property in the request string to the path passed as the JSON `destination`
// property.
func WebMove(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	params, err := ParseRequestBodyParams(r)
	if err != nil {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	source, ok := params["source"]
	if !ok {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	destination, ok := params["destination"]
	if !ok {
		return nil, &augeas.Error{Code: augeas.EBADARG}
	}
	err = wrapper.Move(source, destination)
	return ObjectMap{}, err
}

// WebReload reloads the Augeas configuration tree from the mapped configuration
// files.
func WebReload(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	err := wrapper.Load()
	return ObjectMap{}, err
}

// WebSave persists the Augeas configuration tree into the mapped configuration
// files.
func WebSave(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	err := wrapper.Save()
	return ObjectMap{}, err
}
