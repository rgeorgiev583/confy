// api.go
package wrapper

import (
	"encoding/json"
	"net/http"
	"strings"

	"honnef.co/go/augeas"
)

type ConfigWrapper struct {
	augeas.Augeas
}

type SerializableError struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	MinorMessage string `json:"minor_message"`
	Details      string `json:"details"`
}

type ParamMap map[string]string
type ObjectMap map[string]interface{}

type RequestHandler func(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error)

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

func ParseRequestBodyParams(r *http.Request) (ParamMap, error) {
	var body []byte
	r.Body.Read(body)
	params := ParamMap{}
	err := json.Unmarshal(body, params)
	return params, err
}

func (wrapper ConfigWrapper) List(path string) ([]string, error) {
	return wrapper.Match(path + "/*")
}

func WebList(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	return wrapper.List(path)
}

func WebMatch(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	return wrapper.Match(path)
}

func WebGet(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	value, err := wrapper.Get(path)
	return ObjectMap{"value": value}, err
}

func WebGetAll(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	return wrapper.GetAll(path)
}

func WebGetLabel(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	return wrapper.Label(path)
	value, err := wrapper.Label(path)
	return ObjectMap{"label": value}, err
}

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

func WebClear(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	err := wrapper.Clear(path)
	return ObjectMap{}, err
}

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

func WebRemove(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	num := wrapper.Remove(path)
	return ObjectMap{"count": num}, nil
}

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

func WebReload(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	err := wrapper.Load()
	return ObjectMap{}, err
}

func WebSave(wrapper ConfigWrapper, path string, r *http.Request) (interface{}, error) {
	err := wrapper.Save()
	return ObjectMap{}, err
}
