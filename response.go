package knockttp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	http.Response
	*bytes.Buffer
}

func (m *Response) Header() http.Header {
	return m.Response.Header
}

func (m *Response) Write(p []byte) (int, error) {
	i, err := m.Buffer.Write(p)
	m.Response.ContentLength = m.Response.ContentLength + int64(i)
	return i, err
}

func (m *Response) WriteHeader(statusCode int) {
	m.Response.StatusCode = statusCode
	m.Response.Status = fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
}

func NewResponse(r *http.Request) *Response {
	buffer := bytes.NewBuffer([]byte{})

	w := Response{
		Buffer: buffer,
		Response: http.Response{
			Proto:      r.Proto,
			ProtoMajor: r.ProtoMajor,
			ProtoMinor: r.ProtoMinor,
			Header:     make(http.Header),
			Request:    r,
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(buffer),
		},
	}
	return &w
}
