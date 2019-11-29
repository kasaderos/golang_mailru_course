package middleware

import (
	"fmt"
	"net/http"

	"redditclone/pkg/session"
)

var (
	noAuthUrls = map[string]struct{}{
		"/api/login":    struct{}{},
		"/api/register": struct{}{},
	}
	noSessUrls = map[string]struct{}{
		"/":           struct{}{},
		"/api/posts/": struct{}{},
		"/static/css/main.74225161.chunk.css.map": struct{}{},
		"/static/js/main.32ebaf54.chunk.js.map":   struct{}{},
		"/static/js/2.d59deea0.chunk.js.map":      struct{}{},
		"/static/js/main.32ebaf54.chunk.js":       struct{}{},
		"/static/js/2.d59deea0.chunk.js":          struct{}{},
		"/static/css/main.74225161.chunk.css":     struct{}{},
	}
)

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware")
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		//sess, err := sm.Check(r)

		//_, canbeWithouthSess := noSessUrls[r.URL.Path]
		//fmt.Println("CANBE: ", canbeWithouthSess, "ERROR:", err.Error())
		//if err != nil && !canbeWithouthSess {
		//	fmt.Println("no auth")
		//	http.Redirect(w, r, "/", 302)
		//	return
		//}
		//ctx := context.WithValue(r.Context(), session.SessionKey, sess)
		//next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r)
	})
}
