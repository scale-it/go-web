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
//    curl -i http://0.0.0.0:8000/log and http://0.0.0.0:8000/log/other to see logs
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	goweb "github.com/scale-it/go-web"
	"github.com/scale-it/go-web/contentnegotiator"
	"github.com/scale-it/go-web/handlers"
)

var counter int

func LogHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Running in log handler.")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, int) {
	counter += 1
	return "simple.html",
		fmt.Sprintln("Hello, world! counter=", counter), 200
}

func DataHandler(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	counter += 1
	return counter, 200
}

type Logger struct {
	l *log.Logger
}

func (l Logger) Error(v ...interface{}) {
	l.l.Fatal(v...)
}

func main() {
	var logger1 = log.New(os.Stderr, "[main] ", log.LstdFlags)
	var logger2 = Logger{log.New(os.Stderr, "[log-test] ", log.LstdFlags)}
	t := template.Must(template.ParseGlob("../templates/*.html"))

	// here we use XHandler, which preserve status code and support logger function.
	http.Handle("/log", handlers.XHandler{
		Logger: func(req *http.Request, path string, created time.Time, status, bytes int) {
			goweb.LogRequest(func(s string) { logger1.Print(s) },
				req, created, status, bytes)
		},
		Handler: http.HandlerFunc(LogHandler)})
	http.Handle("/", contentnegotiator.TRenderer{logger2, t, IndexHandler})
	http.Handle("/data", contentnegotiator.Renderer{logger2, DataHandler})
	logger1.Println("Starting listening ...")
	http.ListenAndServe(":8000", nil)
}
