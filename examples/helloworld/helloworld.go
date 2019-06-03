// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Same as the default net/http.

package main

import (
	"fmt"
	"net/http"
	"time"

	goweb "github.com/scale-it/go-web"
	"github.com/scale-it/go-web/handlers"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world")
}

func main() {
	http.Handle("/", withRequestLogger(http.HandlerFunc(indexHandler)))
	srv := http.Server{
		Addr: ":8080",
	}
	srv.ListenAndServe()
}

func logRequest(s string) { fmt.Println(s) }

// WithRequestLogger tracks the request HTTP info and time
func withRequestLogger(handler http.Handler) http.Handler {
	return handlers.XHandler{
		Logger: func(r *http.Request, path string, created time.Time, status, bytes int) {
			goweb.LogRequest(logRequest, r, created, status, bytes)
		},
		Handler:  handler,
		XHeaders: true,
	}
}
