package knockttp

import (
	"io/ioutil"
	"net/http"
	"text/template"
)

// WildcardMethod points to the default handler when there is no explicit method defined
var WildcardMethod = "*"

//Handler contains the elements to build a response
type Handler struct {
	ContentType string `json:"content_type,omitempty"`
	StatusCode  int    `json:"status_code,omitempty"`
	//Location is specified if a redirect is required
	Location string      `json:"location,omitempty"`
	Filename string      `json:"filename,omitempty"`
	Body     []byte      `json:"body,omitempty"`
	Headers  http.Header `json:"headers,omitempty"`
	template *template.Template
}

//Handle processes a http response
func (m *Handler) Handle(w http.ResponseWriter, r *http.Request, Data Data) {
	statusCode := m.StatusCode
	if m.Location != "" {
		if statusCode == 0 {
			statusCode = http.StatusMovedPermanently
		}
		http.Redirect(w, r, m.Location, statusCode)
	}

	contentType := m.ContentType
	if contentType == "" {
		contentType = "text/html"
	}
	w.Header().Set("Content-Type", contentType)

	if m.Headers != nil {
		for k, h := range m.Headers {
			for _, v := range h {
				w.Header().Add(k, v)
			}
		}
	}

	if m.Filename != "" {
		buff, err := ioutil.ReadFile(m.Filename)
		if err != nil {
			panic(err)
		}
		if m.template == nil {
			tmpl := template.New("")
			if m.template, err = tmpl.Parse(string(buff)); err != nil {
				panic(err)
			}
		}
		if err := m.template.Execute(w, Data); err != nil {
			statusCode = http.StatusInternalServerError
			w.Write([]byte(err.Error()))
		}
	} else if string(m.Body) != "" {
		w.Write(m.Body)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(http.StatusText(statusCode)))
	}

	if statusCode == 0 {
		statusCode = 200
	}
	w.WriteHeader(statusCode)
}

//MethodHandlers contains handlers for specific methods
type MethodHandlers map[string]*Handler

//GetHandler returns a Handler given a Method
func (m *MethodHandlers) GetHandler(Method string) (*Handler, bool) {
	h, ok := (*m)[Method]
	if ok {
		return h, ok
	}
	h, ok = (*m)[WildcardMethod]
	return h, ok
}
