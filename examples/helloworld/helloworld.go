// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Same as the default net/http.

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/scale-it/go-web/httpxtra"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world")
}

func main() {
	http.HandleFunc("/", IndexHandler)
	srv := http.Server{
		Addr:    ":8080",
		Handler: httpxtra.Handler{Logger: logger},
	}
	srv.ListenAndServe()
}

func logger(r *http.Request, created time.Time, status, bytes int) {
	fmt.Println(httpxtra.ApacheCommonLog(r, created, status, bytes))
}
