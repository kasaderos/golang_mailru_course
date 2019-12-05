package middleware

import (
	"context"
	"fmt"
	"net/http"
	"redditclone/pkg/session"
	"strings"
)

var (
	noAuthUrls = map[string]struct{}{
		"/api/login":    struct{}{},
		"/api/register": struct{}{},
	}
	noSessUrls = map[string]struct{}{
		"/":                                   struct{}{},
		"/static/js/main.32ebaf54.chunk.js":   struct{}{},
		"/static/js/2.d59deea0.chunk.js":      struct{}{},
		"/static/css/main.74225161.chunk.css": struct{}{},
	}
)

func canbeWithouthSess(r *http.Request) bool {
	for url, _ := range noSessUrls {
		if strings.HasPrefix(r.URL.Path, url) {
			return true
		}
	}
	return false
}

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware")
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)
		if err != nil && (!canbeWithouthSess(r) || r.Method != "GET") {
			fmt.Println("no auth")
			http.Redirect(w, r, "/", 302)
			return
		}
		fmt.Println("session", sess)
		ctx := context.WithValue(r.Context(), session.SessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
