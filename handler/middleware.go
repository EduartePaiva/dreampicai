package handler

import (
	"context"
	"dreampicai/pkg/sb"
	"dreampicai/types"
	"fmt"
	"net/http"
	"strings"
)

// type contextKey string

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}
		user := types.AuthenticatedUser{}
		cookie, err := r.Cookie("at")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		resp, err := sb.Client.Auth.User(r.Context(), cookie.Value)
		if err == nil {
			user.Email = resp.Email
			user.LoggedIn = true
			fmt.Println("foi aqui")
		} else {
			fmt.Println("foi aqui2")
			fmt.Println(err)
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)

}
