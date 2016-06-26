package gas

import (
	// "bufio"
	// "fmt"
	// "net"
	"net/http"
)

type ResponseWriter struct {
	Writer http.ResponseWriter
	status int
	size   int
}

func (rw *ResponseWriter) reset(w http.ResponseWriter) {
	rw.Writer = w
	rw.size = 0
	rw.status = http.StatusOK
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.Writer.Header()
}

func (rw *ResponseWriter) WriteHeader(s int) {
	rw.status = s
	rw.Writer.WriteHeader(s)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.Writer.Write(b)
	rw.size += size
	return size, err
}

func (rw *ResponseWriter) Status() int {
	return rw.status
}

func (rw *ResponseWriter) Size() int {
	return rw.size
}

func (rw *ResponseWriter) Written() bool {
	return rw.status != 0
}

func (rw *ResponseWriter) CloseNotify() <-chan bool {
	return rw.Writer.(http.CloseNotifier).CloseNotify()
}

func (rw *ResponseWriter) Flush() {
	flusher, ok := rw.Writer.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
