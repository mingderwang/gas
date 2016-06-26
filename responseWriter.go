package gas

import (
	// "bufio"
	// "fmt"
	// "net"
	"net/http"
	"github.com/valyala/fasthttp"
)

type ResponseWriter struct {
	Writer fasthttp.Response
	status int
	size   int
}

func (rw *ResponseWriter) reset(w fasthttp.Response) {
	rw.Writer = w
	rw.size = 0
	rw.status = http.StatusOK
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.Writer.Header()
}

func (rw *ResponseWriter) WriteHeader(s int) {
	rw.status = s
	rw.Writer.Header().SetStatusCode(s)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	err := rw.Writer.Write(b)
	rw.size += rw.Writer
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

//func (rw *ResponseWriter) CloseNotify() <-chan bool {
//	return rw.Writer.(http.CloseNotifier).CloseNotify()
//}
//
//func (rw *ResponseWriter) Flush() {
//	flusher, ok := rw.Writer.(http.Flusher)
//	if ok {
//		flusher.Flush()
//	}
//}
