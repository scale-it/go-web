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

// Run this, and from console:
//    curl -i -H "Accept: application/msgpack"  http://localhost:8000/data
//    curl -i -H "Accept: application/json"  http://localhost:8000/data
//    curl -i http://localhost:8000/data
//    curl -i http://localhost:8000
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/scale-it/go-web/contentnegotiator"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, int) {
	counter += 1
	return "simple.html", counter, 200
}

func DataHandler(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	counter += 1
	return counter, 200
}

var counter int
var logger = log.New(os.Stderr, "renderer", log.LstdFlags)

func main() {
	t := template.Must(template.ParseGlob("../templat es/*.html"))
	http.Handle("/", contentnegotiator.TRenderer{logger, t, IndexHandler})
	http.Handle("/data", contentnegotiator.Renderer{logger, DataHandler})
	http.ListenAndServe(":8000", nil)
}
