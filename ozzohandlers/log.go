// Copyright (c) 2017 Robert Zaremba
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ozzohandlers

import (
	"time"

	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/robert-zaremba/log15"
	goweb "github.com/scale-it/go-web"
	"github.com/scale-it/go-web/handlers"
)

// LogTrace is a structure which provides ozzo routre.Handler
// to track the request HTTP info and time. It also logs the server errors.
type LogTrace struct {
	Logger log15.Logger
}

func (lt LogTrace) logRequest(s string) { lt.Logger.Debug(s) }

// LogTrace implements ozzo routing.Handler interface
func (lt LogTrace) LogTrace(c *routing.Context) error {
	t := time.Now()
	w := handlers.WrapWriter(c.Response)
	c.Response = w
	// originalPath := r.URL.Path // this can be overwritten by a middleware
	err := c.Next()
	status := w.Status()
	if err != nil {
		if errRouting, ok := err.(routing.HTTPError); ok {
			status = errRouting.StatusCode()
		}
	} else if status == 0 {
		// We overwrite status only locally for the Trace
		status = 200
	}
	goweb.LogRequest(lt.logRequest, c.Request, t, status, w.BytesWritten())
	return err
}
