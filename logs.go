package goweb

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

// Function which construct Apache like message with information
// about request duration
func LogRequest(log func(string), req *http.Request, created time.Time, status, bytes int) {
	username := "-"
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}
	elapsed := float64(time.Since(created)) / float64(time.Millisecond)
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		ip = req.RemoteAddr
	}

	log(fmt.Sprintf("%s - %s \"%s %s %s\" %d %dB \"%s\". %fms",
		ip,
		username,
		req.Method,
		req.RequestURI,
		req.Proto,
		status,
		bytes,
		req.UserAgent(),
		elapsed))
}
