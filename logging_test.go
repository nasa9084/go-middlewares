package middlewares_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	middlewares "github.com/nasa9084/go-middlewares"
)

func TestLogger(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetOutput(&buf)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(`GET`, `/`, nil)
	if err != nil {
		panic(err)
	}

	mwset := middlewares.New(middlewares.Logger)
	mwset.ApplyFunc(testHandler).ServeHTTP(w, r)

	logline := strings.Split(buf.String(), "\n")

	// request log test
	requestLog := strings.SplitN(logline[0], " ", 3)[2]
	expectRequestLog := "GET /"
	if requestLog != expectRequestLog {
		t.Errorf(`"%s" != "%s"`, requestLog, expectRequestLog)
		return
	}

	// response log test
	responseLog := strings.SplitN(logline[1], " ", 3)[2]
	expectResponseLog := "200 OK"
	if responseLog != expectResponseLog {
		t.Errorf(`"%s" != "%s"`, responseLog, expectResponseLog)
		return
	}
}
