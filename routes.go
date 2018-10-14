package knockttp

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Route struct {
	Path    string         `json:"query,omitempty"`
	Methods MethodHandlers `json:"methods,omitempty"`
}

func (r *Route) GetHandler(Method string) (*Handler, bool) {
	return r.Methods.GetHandler(Method)
}

type Routes []*Route

func NewRoutesFromFile(Filename string) (*Routes, error) {
	f, err := os.Open(Filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewRoutesFromReader(f)
}

func NewRoutesFromReader(r io.Reader) (*Routes, error) {
	var routes Routes
	dec := json.NewDecoder(r)
	if err := dec.Decode(&routes); err != nil {
		return nil, err
	}
	return &routes, nil
}

func (q *Routes) Find(r *http.Request) (*Route, bool) {
	for _, route := range *q {
		if route.Path == r.URL.Path {
			return route, true
		}
	}
	return nil, false
}

func (r *Routes) GetHandler(req *http.Request) (*Route, *Handler, bool) {
	if route, ok := r.Find(req); ok {
		if handler, ok := route.GetHandler(req.Method); ok {
			return route, handler, true
		}
	}
	return nil, nil, false
}
