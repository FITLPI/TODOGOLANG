package core_http_middleware

import "net/http"

type Middleware func(next http.Handler) http.Handler

func ChainMidlewares(h http.Handler, m ...Middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
