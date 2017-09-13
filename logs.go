package goweb

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// LogRequest constructs Apache like log with request duration
func LogRequest(log func(string), req *http.Request, created time.Time, status, bytes int) {
	username := "-"
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}
	elapsed := float64(time.Since(created)) / float64(time.Millisecond)
	ip := GetClientIP(req)

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

// GetClientIP retrives request client IP
func GetClientIP(req *http.Request) string {
	ip := req.Header.Get("X-Real-IP")
	if ip == "" {
		ip = req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = req.RemoteAddr
		}
	}
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}
