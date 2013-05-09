/* Copyright 2013 Robert Zaremba
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"net/http"
	"strings"
)

// Handler which check if request is from canonical host and uses https. Otherwise will
// redirect to https://<canonicalhost>/rest/of/the/url
type ForceHTTPS struct {
	CanonicalHost string
	Next          http.Handler
}

func (this ForceHTTPS) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	is_http := true
	if h, ok := req.Header["X-Forwarded-Proto"]; ok {
		if h[0] == "https" {
			is_http = false
		}
	}
	hostPort := strings.Split(req.Host, ":")
	if is_http || hostPort[0] != this.CanonicalHost {
		hostPort[0] = this.CanonicalHost
		url := "https://" + strings.Join(hostPort, ":") + req.URL.String()
		http.Redirect(w, req, url, http.StatusMovedPermanently)
		return
	}

	this.Next.ServeHTTP(w, req)
}

// Ensures authentication for handlers. Otherwise call fallback.
type Auth struct {
	A        Authenticator
	Next     http.Handler
	Fallback http.Handler
}

type Authenticator func(req *http.Request) bool

func (this Auth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if this.A(req) {
		this.Next.ServeHTTP(w, req)
	} else {
		this.Fallback.ServeHTTP(w, req)
	}
}

// Calls the wrapped handler and on panic calls the specified error handler.
// errH can make some logging or just return:
//   http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
func PanicHandler(h http.Handler, errH func(http.ResponseWriter, *http.Request, interface{})) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errH(w, r, err)
			}
		}()
		h.ServeHTTP(w, r)
	}
}

// Handler which chains other multiple handlers into single one
func Chain(h ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range h {
			v.ServeHTTP(w, r)
		}
	})
}

// Helper function to handle http.HandlerFunc
func ChainFuncs(h ...http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range h {
			v(w, r)
		}
	})
}
