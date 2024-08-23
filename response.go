package pico

import (
	"bytes"
	"fmt"
	"net/http"
)

func newHttpResponseWriter() *httpResponseWriter {
	return &httpResponseWriter{
		buffer: bytes.NewBuffer(nil),
		header: make(http.Header),
	}
}

type httpResponseWriter struct {
	buffer *bytes.Buffer
	header http.Header

	wroteHeader bool
}

func (rw *httpResponseWriter) Header() http.Header {
	return rw.header
}

func (rw *httpResponseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.buffer.Write(b)
}

func (rw *httpResponseWriter) WriteHeader(statusCode int) {
	ln := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))
	rw.buffer.Write([]byte(ln))
	_ = rw.header.Write(rw.buffer)
	rw.wroteHeader = true
}

func (rw *httpResponseWriter) Bytes() []byte {
	return rw.buffer.Bytes()
}
