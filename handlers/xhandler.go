package handlers

import (
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
	wp := WrapWriter(w)
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
	originalPath := r.URL.Path // this can be overwritten by a middleware
	h.Handler.ServeHTTP(wp, r)
	if h.Logger != nil {
		h.Logger(r, originalPath, t, wp.Status(), wp.BytesWritten())
	}
}

// LoggerFunc can be called by XHandler at the end of each request.
type LoggerFunc func(r *http.Request, patch string, created time.Time, status, bytes int)
