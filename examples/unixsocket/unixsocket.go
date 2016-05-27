// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// unixsocket.go starts an HTTP server on a Unix socket instead of TCP port.
//
// This is useful when the server is reverse proxied by a frontend server or
// a load balancer like Nginx.
//
// Make sure the frontend server sets either X-Real-IP or X-Forwarded-For HTTP
// headers with the IP address of the client, and set XHeaders=true in our
// custom httpxtra.Handler below.
//
// When XHeaders is set to true, it overwrites http.Request.RemoteAddr with
// the contents of either X-Real-IP or X-Forwarded-For HTTP header. The IP is
// not validated.
//
// Test the server:
// echo -ne 'GET / HTTP/1.1\r\nX-Real-IP: pwnz\r\n\r\n' | nc -U ./test.sock
package main

import (
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/scale-it/go-web/httpxtra"
)

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello, world")
}

func main() {
	http.HandleFunc("/", IndexHandler)

	// Setup the custom handler
	handler := httpxtra.Handler{
		Logger:   logger,
		XHeaders: true,
	}

	// Setup the server
	server := http.Server{
		Addr:    "./test.sock", // Listen on Unix Socket
		Handler: handler,       // Custom httpxtra.Handler
	}

	// ListenAndServe fails with "address already in use" if the socket
	// file exists.
	syscall.Unlink("./test.sock")

	// Use our custom listener
	if e := httpxtra.ListenAndServe(server); e != nil {
		fmt.Println(e.Error())
	}
}

func logger(r *http.Request, created time.Time, status, bytes int) {
	fmt.Println(httpxtra.ApacheCommonLog(r, created, status, bytes))
}
