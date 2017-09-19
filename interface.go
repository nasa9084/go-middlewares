package middlewares

import "net/http"

// Middleware is a http handling middleware
type Middleware func(http.Handler) http.Handler
