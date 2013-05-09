package handlers

import (
	"bufio"
	"net"
	"net/http"
	"time"
)

// XtraHandler is wrapper for http.Handler that adds extra features to the server:
// - Custom logging
// - Support for listening on TCP or UNIX sockets
// - Support X-Real-IP and X-Forwarded-For as the remote IP if the server sits
//   behind a proxy or load balancer.
type XHandler struct {
	Handler  http.Handler
	Logger   LoggerFunc
	XHeaders bool
}

func (h XHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	lw := logWriter{w: w}
	if h.Handler == nil {
		h.Handler = http.DefaultServeMux
	}
	if h.XHeaders {
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = r.Header.Get("X-Forwarded-For")
		}
		if ip != "" {
			r.RemoteAddr = ip
		}
	}
	h.Handler.ServeHTTP(&lw, r)
	if h.Logger != nil {
		h.Logger(r, t, lw.status, lw.bytes)
	}
}

// LoggerFunc can be called by XHandler at the end of each request.
type LoggerFunc func(r *http.Request, created time.Time, status, bytes int)

type logWriter struct {
	w      http.ResponseWriter
	bytes  int
	status int
}

func (lw *logWriter) Header() http.Header {
	return lw.w.Header()
}

func (lw *logWriter) Write(b []byte) (int, error) {
	if lw.status == 0 {
		lw.status = http.StatusOK
	}
	n, err := lw.w.Write(b)
	lw.bytes += n
	return n, err
}

func (lw *logWriter) WriteHeader(s int) {
	lw.w.WriteHeader(s)
	lw.status = s
}

func (lw *logWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if lw.status == 0 {
		lw.status = http.StatusOK
	}
	// TODO: Check. Does it break if the server don't support hijacking?
	return lw.w.(http.Hijacker).Hijack()
}
