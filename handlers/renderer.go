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
	"encoding/json"
	"fmt"
	"github.com/hinasssan/msgpack-go"
	"github.com/scale-it/go-log"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const (
	r_error = iota
	r_json
	r_msgpack
	r_xml
	r_unknown
)

// Structure renderer. It renders the handler output using encoders (json, msgpack ...).
// The encoder is chose by request "Content-type" header
type Renderer struct {
	Log *log.Logger
	/* handler, which output will be rendered. It should return
	 * data to be rendered. data is an error, then http.Error will be used to render it.
	 * status code */
	H func(w http.ResponseWriter, r *http.Request) (interface{}, int)
}

func (this Renderer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, status := this.H(w, r)
	if data_error, ok := data.(error); ok {
		http.Error(w, data_error.Error(), status)
	}
	switch negotiateRenderer(r.Header.Get("Accept")) {
	case r_json:
		w.Header().Set("Content-Type", "application/json")
		content, err := json.Marshal(data)
		write(this.Log, w, content, err, status)
	case r_msgpack:
		w.Header().Set("Content-Type", "application/x-msgpack")
		content, err := msgpack.Marshal(data)
		write(this.Log, w, content, err, status)
	default:
		w.Header().Set("Content-Type", "text/plain")
		write(this.Log, w, []byte(fmt.Sprint(data)), nil, status)
	}
}

// Template renderer. It renders the handler output using http.template.
type TRenderer struct {
	Log *log.Logger
	T   *template.Template
	/* handler, which output will be rendered. It should return
	* template name which is a fielname associated to `T`.
	* data to be rendered. data is an error, then http.Error will be used to render it.
	* status code */
	H func(w http.ResponseWriter, r *http.Request) (string, interface{}, int)
}

func (this TRenderer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tname, data, status := this.H(w, r)
	if data_error, ok := data.(error); ok {
		http.Error(w, data_error.Error(), status)
	}
	w.Header().Set("Content-Type", "text/html")
	if err := this.T.ExecuteTemplate(w, tname, data); err != nil {
		write(this.Log, w, nil, err, status)
	}
}

func write(logger *log.Logger, w http.ResponseWriter, data []byte, err error, status int) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error(err.Error())
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

func negotiateRenderer(field string) int {
	for _, a := range strings.Split(field, ",") {
		if strings.Contains(a, "json") {
			return r_json
		}
		if strings.Contains(a, "msgpack") {
			return r_msgpack
		}
		if strings.Contains(a, "xml") {
			return r_xml
		}
	}
	return r_unknown
}
