package middlewares_test

import (
	"fmt"
	"net/http"
	"testing"

	middlewares "github.com/nasa9084/go-middlewares"
)

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world\n"))
})

func equalFunc(fn1, fn2 interface{}) bool {
	return fmt.Sprintf("%p", fn1) == fmt.Sprintf("%p", fn2)
}

func TestNew(t *testing.T) {
	mws := []middlewares.Middleware{
		func(h http.Handler) http.Handler {
			return nil
		},
	}
	mwset := middlewares.New(mws...)
	for i := range mws {
		if !equalFunc(mwset[i], mws[i]) {
			t.Errorf("%v != %v", &mwset[i], &mws[i])
			return
		}
	}
}

func TestApply(t *testing.T) {
	candidates := []struct {
		name string
		mwset  middlewares.Middlewareset
		target http.Handler
		expect http.Handler
		msg    string
	}{
		{"empty mwset", middlewares.New(), testHandler, testHandler, "empty middlewareset should not affect"},
		{"nil handler", middlewares.New(), nil, http.DefaultServeMux, "apply to nil should return http.DefaultServeMux"},
	}
	for _, c := range candidates {
		t.Log(c.name)
		if !equalFunc(c.mwset.Apply(c.target), c.expect) {
			t.Errorf("%s != %s", c.mwset.Apply(c.target), c.expect)
			t.Errorf(c.msg)
			return
		}
	}
}
