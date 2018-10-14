package knockttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"text/template"
)

type TestServer struct {
	*httptest.Server

	Data   map[string]interface{}
	Routes Routes
}

func NewTestServer(routes Routes) (*TestServer, error) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	ts := TestServer{
		Server: server,
		Data: map[string]interface{}{
			"BaseURL": server.URL,
		},
		Routes: routes,
	}

	for _, route := range ts.Routes {
		for method, handler := range route.Methods {
			if handler.Filename != "" {
				buff, err := ioutil.ReadFile(handler.Filename)
				if err != nil {
					panic(err)
				}
				t := template.New(method + " " + route.Path)
				handler.template, err = t.Parse(string(buff))
				if err != nil {
					return nil, err
				}
			}
		}
		mux.HandleFunc(route.Path, ts.ServeFunc(route))
	}

	return &ts, nil
}

func (ts *TestServer) ServeFunc(route *Route) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		m, ok := route.GetHandler(r.Method)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(fmt.Sprintf("unsupported method: '%s'", r.Method)))
			return
		}
		m.Handle(w, r, ts.Data)
	}
}
