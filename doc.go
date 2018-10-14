/*
Package knockttp is a Go library providing to mock a http server for testing purposes

It has been initially designed to test crawlers reproducing different web scenarios.

There are two modes of operation:

Transport - works by setting a http.RoundTripper in your client requests. Due
the nature of this implementation, a multiple host/service setup can be implemented.

Server - creates a test server into a random port

*/
package knockttp
