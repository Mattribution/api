package http

import (
	"log"
	"net/http"
)

// Define our struct
type authMiddleware struct {
	tokenUsers map[string]int
	h          Handler
}

func newAuthMiddleware(h Handler) authMiddleware {
	amw := authMiddleware{
		h: h,
	}
	amw.Populate()
	return amw
}

// Initialize it somewhere
func (amw *authMiddleware) Populate() {
	amw.tokenUsers["00000000"] = 1
}

// Middleware function, which will be called for each request
func (amw *authMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)

			// TODO:
			// Find the user's data sources
			// if contains presto replacement
			//   replace presto service
			//
			// ENDTODO:

			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
