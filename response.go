package pico

import (
	"bytes"
	"fmt"
	"net/http"
)

func NewHttpResponseWriter() *HttpResponseWriter {
	return &HttpResponseWriter{
		buffer: bytes.NewBuffer(nil),
		header: make(http.Header),
	}
}

type HttpResponseWriter struct {
	buffer *bytes.Buffer
	header http.Header

	wroteHeader bool
}

func (rw *HttpResponseWriter) Header() http.Header {
	return rw.header
}

func (rw *HttpResponseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.buffer.Write(b)
}

func (rw *HttpResponseWriter) WriteHeader(statusCode int) {
	ln := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))
	rw.buffer.Write([]byte(ln))
	_ = rw.header.Write(rw.buffer)
	rw.wroteHeader = true
}

func (rw *HttpResponseWriter) Bytes() []byte {
	return rw.buffer.Bytes()
}
