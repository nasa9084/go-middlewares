package middlewares

import (
	"log"
	"net/http"
	"strings"
)

type logResponseWriter struct {
	http.ResponseWriter
	status   int
	respBody []byte
}

func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *logResponseWriter) Header() http.Header {
	return lrw.ResponseWriter.Header()
}

func (lrw *logResponseWriter) Write(b []byte) (int, error) {
	lrw.respBody = b
	return lrw.ResponseWriter.Write(b)
}

// Logger logs request with log.DefaultLogger
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := logResponseWriter{
			ResponseWriter: w,
			status:         200,
		}
		log.Printf("%s %s", r.Method, r.URL.RequestURI())
		h.ServeHTTP(&lrw, r)
		log.Printf("  => %d %s", lrw.status, http.StatusText(lrw.status))
		if lrw.respBody != nil {
			for _, line := range strings.Split(string(lrw.respBody), "\n") {
				log.Printf("  > %s", line)
			}
		}
	})
}
