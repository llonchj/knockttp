package knockttp

import (
	"fmt"
	"net/http"
)

type Data map[string]interface{}

type Transport struct {
	Network
	Data
}

//NewTransport instantiates a transport given a Network and Data
func NewTransport(Network Network, Data Data) *Transport {
	return &Transport{Network: Network, Data: Data}
}

//RoundTrip implements http.RoundTripper
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	host := r.Host
	if host == "" {
		host = r.URL.Hostname()
	}
	routes, ok := (*t).Network.Services[Host(host)]
	if !ok {
		return nil, fmt.Errorf("not found: '%s'", host)
	}

	w := NewResponse(r)
	_, method, ok := routes.GetHandler(r)
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return &w.Response, nil
	}
	method.Handle(w, r, t.Data)
	return &w.Response, nil
}
