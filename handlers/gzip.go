// Copyright 2013 Alexandre Fiori, Robert Zaremba
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE.bsd file.

// autogzip provides on-the-fly gzip encoding for http servers.
package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type IOResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w IOResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Handle provides on-the-fly gzip encoding for other handlers.
//
// Usage:
//
//	func DL1Handler(w http.ResponseWriter, req *http.Request) {
//		fmt.Fprintln(w, "foobar")
//	}
//
//	func DL2Handler(w http.ResponseWriter, req *http.Request) {
//		fmt.Fprintln(w, "zzz")
//	}
//
//
//	func main() {
//		http.HandleFunc("/download1", DL1Handler)
//		http.HandleFunc("/download2", DL2Handler)
//		http.ListenAndServe(":8080", autogzip.Handle(http.DefaultServeMux))
//	}
func Gzip(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Do nothing on a HEAD request
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || r.Method == "HEAD" {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		h.ServeHTTP(IOResponseWriter{Writer: gz, ResponseWriter: w}, r)
	}
}
