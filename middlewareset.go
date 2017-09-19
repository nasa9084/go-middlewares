package middlewares

import "net/http"

// Middlewareset is a set of Middlewares
type Middlewareset []Middleware

// New returns a new Middlewareset
func New(mws ...Middleware) Middlewareset {
	return mws
}

// NewSubgroup creates new sub-Middlewareset
func (mwset Middlewareset) NewSubgroup(mws ...Middleware) Middlewareset {
	return append(mwset, mws...)
}

// Extend is an alias to NewSubgroup
func (mwset Middlewareset) Extend(mws ...Middleware) Middlewareset {
	return mwset.NewSubgroup(mws...)
}

// Apply middlewares to handler
func (mwset Middlewareset) Apply(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}
	for _, mw := range mwset {
		h = mw(h)
	}
	return h
}

// ApplyFunc applies middlewares to HandlerFunc
func (mwset Middlewareset) ApplyFunc(hfn http.HandlerFunc) http.Handler {
	return mwset.Apply(hfn)
}
