package middleware

import (
	"net/http"
	"strings"
)

func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/public/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		}
		next.ServeHTTP(w, r)
	})
}
