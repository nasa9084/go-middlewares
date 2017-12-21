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
	for i := len(mwset) - 1; i >= 0; i-- {
		h = mwset[i](h)
	}
	return h
}

// ApplyFunc applies middlewares to HandlerFunc
func (mwset Set) ApplyFunc(hfn http.HandlerFunc) http.Handler {
	return mwset.Apply(hfn)
}
