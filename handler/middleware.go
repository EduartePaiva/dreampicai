package handler

import (
	"context"
	"dreampicai/pkg/sb"
	"dreampicai/types"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
)

// type contextKey string

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
		session, err := store.Get(r, sessionUserKey)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		accessToken := session.Values[sessionAccessTokenKey]

		resp, err := sb.Client.Auth.User(r.Context(), accessToken.(string))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, types.AuthenticatedUser{
			Email:    resp.Email,
			LoggedIn: true,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)

}
func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := getAuthenticatedUser(r)
		if !user.LoggedIn {
			path := r.URL.Path
			http.Redirect(w, r, "/login?to="+path, http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)

}
