package middlewares_test

import (
	"log"
	"net/http"

	middlewares "github.com/nasa9084/go-middlewares"
)

func ExampleMiddlewareset() {
	logMw := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Print("some log")
			h.ServeHTTP(w, r)
		})
	}
	mwset := middlewares.New(logMw)

	if err := http.ListenAndServe(":8080", mwset.Apply(nil)); err != nil {
		log.Printf("%s", err)
		return
	}
}
