// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//
// This demo is live at http://cos.pe

package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/scale-it/go-web/sse"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	sf := 0
	startFrame := r.FormValue("startFrame")
	if startFrame != "" {
		sf, _ = strconv.Atoi(startFrame)
	}
	if sf < 0 || sf >= cap(frames) {
		sf = 0
	}
	conn, buf, err := sse.ServeEvents(w)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	// Play the movie, frame by frame
	for n, f := range frames[sf:] {
		m := &sse.MessageEvent{Id: strconv.Itoa(n + 1), Data: f.Buf}
		e := sse.SendEvent(buf, m)
		if e != nil {
			// usually a broken pipe error
			// log.Println(e.Error())
			break
		}
		time.Sleep(f.Time)
	}
}

func main() {
	err := loadMovie("./ASCIImation.txt.gz")
	if err != nil {
		log.Println(err)
		return
	}
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/sse", SSEHandler)
	server := http.Server{
		Addr: ":8080",
	}
	server.ListenAndServe()
}

type Message struct {
	FrameNo  int
	FrameBuf string
}

type Frame struct {
	Time time.Duration
	Buf  string // This is a JSON-encoded Message{FrameBuf:...}
}

var frames []Frame

func loadMovie(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	gzfile, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	lineno := 1
	reader := bufio.NewReader(gzfile)
	frameNo := 1
	frameBuf := ""
	var frameTime time.Duration
	var part string
	for {
		if part, err = reader.ReadString('\n'); err != nil {
			break
		}

		switch lineno % 14 {
		case 0:
			b := html.EscapeString(frameBuf + part)
			j, _ := json.Marshal(Message{frameNo, b})
			frames = append(frames, Frame{frameTime, string(j)})
			frameNo++
			frameBuf = ""
		case 1:
			s := string(part)
			n, e := strconv.Atoi(s[:len(s)-1])
			if e == nil {
				frameTime = time.Duration(n) * time.Second / 10
			}
		default:
			frameBuf += part
		}
		lineno += 1
	}
	return nil
}
