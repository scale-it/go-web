// Copyright 2013 Alexandre Fiori, Robert Zaremba
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE.bsd file.

// autogzip provides on-the-fly gzip encoding for http servers.
package handlers

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type IOResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w IOResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
func (w IOResponseWriter) WriteHeader(i int) {
	w.ResponseWriter.WriteHeader(i)
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
//		http.ListenAndServe(":8080", Gzip(http.DefaultServeMux))
//	}
func Gzip(h http.Handler) http.HandlerFunc {
	var pool sync.Pool
	pool.New = func() interface{} {
		return gzip.NewWriter(ioutil.Discard)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Do nothing on a HEAD request
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || r.Method == "HEAD" ||
			w.Header().Get("Content-Encoding") == "gzip" { // Skip compression if already compressed

			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := pool.Get().(*gzip.Writer)
		defer pool.Put(gz)
		gz.Reset(w)

		h.ServeHTTP(IOResponseWriter{Writer: gz, ResponseWriter: WrapWriter(w)}, r)
		gz.Close()
	}
}
