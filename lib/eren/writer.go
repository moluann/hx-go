package eren

import "net/http"

type baseResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher

	// Returns the HTTP response status code of the current request.
	Status() int

	Size() int

	WriteString(string) (int,error)

	Written() bool

	WriteHeaderNow()


}
//custom writer
type responseWriter struct {


}