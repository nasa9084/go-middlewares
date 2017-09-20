package middlewares_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	middlewares "github.com/nasa9084/go-middlewares"
)

const hello = "hello, world"

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(hello))
})

func testHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(hello))
}

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

func TestApplyNil(t *testing.T) {
	candidates := []struct {
		name   string
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

func newPrefixMw(prefix string) middlewares.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(prefix))
			h.ServeHTTP(w, r)
		})
	}
}

func TestApplyOrder(t *testing.T) {
	mwset := middlewares.New(newPrefixMw("1"), newPrefixMw("2"), newPrefixMw("3"))
	applied := mwset.Apply(testHandler)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(`GET`, `/`, nil)
	if err != nil {
		panic(err)
	}

	applied.ServeHTTP(w, r)

	body := w.Body.String()
	if body != "123"+hello {
		t.Errorf(`"%s" != "%s%s"`, body, "123", hello)
		return
	}
}

func TestExtend(t *testing.T) {
	basemwset := middlewares.New(newPrefixMw("1"), newPrefixMw("2"), newPrefixMw("3"))
	candidates := []struct {
		mwset  middlewares.Middlewareset
		expect string
	}{
		{basemwset, "123" + hello},
		{basemwset.Extend(newPrefixMw("4")), "1234" + hello},
	}
	for _, c := range candidates {
		applied := c.mwset.Apply(testHandler)

		w := httptest.NewRecorder()
		r, err := http.NewRequest(`GET`, `/`, nil)
		if err != nil {
			panic(err)
		}

		applied.ServeHTTP(w, r)
		body := w.Body.String()
		if body != c.expect {
			t.Errorf(`"%s" != "%s"`, body, c.expect)
			return
		}
	}
}

func TestApply(t *testing.T) {
	mwset := middlewares.New(newPrefixMw("1"), newPrefixMw("2"), newPrefixMw("3"))
	applied := mwset.Apply(testHandler)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(`GET`, `/`, nil)
	if err != nil {
		panic(err)
	}

	applied.ServeHTTP(w, r)

	body := w.Body.String()
	if body != "123"+hello {
		t.Errorf(`"%s" != "%s"`, body, "123"+hello)
		return
	}
}

func TestApplyFunc(t *testing.T) {
	mwset := middlewares.New(newPrefixMw("1"), newPrefixMw("2"), newPrefixMw("3"))
	applied := mwset.ApplyFunc(testHandlerFunc)
	w := httptest.NewRecorder()
	r, err := http.NewRequest(`GET`, `/`, nil)
	if err != nil {
		panic(err)
	}

	applied.ServeHTTP(w, r)

	body := w.Body.String()
	if body != "123"+hello {
		t.Errorf(`"%s" != "%s"`, body, "123"+hello)
		return
	}
}
