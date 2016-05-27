# go-web

go-web is a collection of utilities for http servers written in Go.

It has the following packages:

- autogzip: An http.Handler that supports on-the-fly gzip encoding.
- httpxtra: Essential logging and support for X-Real-IP and X-Forwarded-For.
- handlers: A set of useful handlers which. Includes httpxtra and gzip functionality.
- remux: A very simple request multiplexer that supports regular expressions.
- sse: Server-Sent Events, a.k.a. HTTP push notifications.

*NOTE*: go-web used to be an experimental fork of Go's
[net/http](http://golang.org/pkg/net/http/) package. It's no longer a fork and
now uses the standard http package, making it easier to be used along with
[Gorilla Toolkit](http://www.gorillatoolkit.org) and other packages.

## Examples

There are some nice [examples](https://github.com/scale-it/go-web/tree/master/examples) including a full featured web application with sign up, using MySQL and
Redis for storage.

*example/handlers* contains step by step introduction to handlers package.

There's also some live stuff:

- freegeoip at http://freegeoip.net
- sse demo at http://cos.pe
