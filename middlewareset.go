package middlewares

import "net/http"

// Set is a set of Middlewares
type Set []Middleware

// New returns a new Middlewareset
func New(mws ...Middleware) Set {
	return mws
}

// NewSubgroup creates new sub-Middlewareset
func (mwset Set) NewSubgroup(mws ...Middleware) Set {
	return append(mwset, mws...)
}

// Extend is an alias to NewSubgroup
func (mwset Set) Extend(mws ...Middleware) Set {
	return mwset.NewSubgroup(mws...)
}

// Apply middlewares to handler
func (mwset Set) Apply(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}
	for _, mw := range reverse(mwset) {
		h = mw(h)
	}
	return h
}

// ApplyFunc applies middlewares to HandlerFunc
func (mwset Set) ApplyFunc(hfn http.HandlerFunc) http.Handler {
	return mwset.Apply(hfn)
}

func reverse(mws []Middleware) []Middleware {
	length := len(mws)
	n := make([]Middleware, length)
	for i := 0; i < length; i++ {
		n[length-i-1] = mws[i]
	}
	return n
}
